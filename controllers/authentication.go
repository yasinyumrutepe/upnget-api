package controllers

import (
	"auction/database"
	"auction/globals"
	"auction/models"
	"auction/secret"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type Authentication struct{}
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterInput struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserData struct {
	SellerName string `json:"seller_name"`
	UserID     uint   `json:"user_id"`
	SellerID   int    `json:"seller_id"`
	Level      uint   `json:"level"`
	ExpiresAt  int64  `json:"ExpiresAt"`
	Token      string `json:"token"`
}

func GetAbilityByUserID(CaslAbilities *models.User) error {
	//DESC - Check user is active. Get employee roles and abilities by user id
	if CaslAbilities.ID != 0 {
		//DESC -Find the user's top role level
		database.Conn.DB.Raw("select role_level from roles join user_roles on roles.id = user_roles.role_id where user_roles.user_id = ? Order by role_level ASC LIMIT 1", CaslAbilities.ID).Row().Scan(&CaslAbilities.Level)
		//DESC - Get user abilities with the module by user_id
		database.Conn.DB.Raw(`SELECT modules.subject as module_subject,abilities.* FROM 
		"abilities" Left join (select DISTINCT user_id,ability_id from user_abilities ) as user_abilities on abilities.id = user_abilities.ability_id  
		left JOIN modules ON modules.id = abilities.module_id WHERE  user_abilities.user_id = ?`, CaslAbilities.ID).Scan(&CaslAbilities.Abilities)
	} else {
		return errors.New("There was an error when bringing user information.Please try again.")
	}
	if CaslAbilities.Level == 0 || CaslAbilities.Abilities == nil {
		return errors.New("You do not have enough authority to enter the system.")
	}

	var tempMapAbility = make(map[string]interface{}, 0)

	//DESC - When users's powers are sent in front, Key Value is sent as Value.
	for _, ability := range CaslAbilities.Abilities {
		if condition, ok := tempMapAbility[ability.ModuleSubject]; ok {
			tempMapAbility[ability.ModuleSubject] = append(condition.([]string), ability.Key)
		} else {
			tempMapAbility[ability.ModuleSubject] = []string{ability.Key}
		}
	}

	//DESC - Set user data to redis
	CaslAbilities.CaslAbilities = tempMapAbility

	userIDStr := strconv.Itoa(int(CaslAbilities.ID))

	err := database.SetName(context.Background(), &CaslAbilities, "User-"+userIDStr, time.
		Duration(1*time.Hour))
	if err != nil {
		fmt.Println("Redis bağlanırken hata oldu")
		fmt.Print(err)
	}
	return nil
}

func HasAbility(ability string, module string, c *fiber.Ctx) (bool, error) {
	userDatas := models.User{}
	userDatas = c.Locals("user").(models.User)
	//DESC - Check user is active. Get employee roles and abilities by user id
	userIDStr := strconv.Itoa(int(userDatas.ID))
	//DESC - Get user data from redis
	datam := database.GetName(context.Background(), "User-"+userIDStr, &userDatas)
	if datam != nil {
		if datam.Error() == "nil" {
			//DESC - User data is not in redis so get user data from db
			err := GetAbilityByUserID(&userDatas)
			if err != nil {
				return false, err
			}
		}
	}

	if condition, ok := userDatas.CaslAbilities[module]; ok {
		if globals.InSlice(ability, condition) {
			return true, nil
		}
	}

	return false, nil

}

func getUserInfo(loginInput *LoginInput) (*models.Seller, error) {
	//DESC - Error Handling
	var err error
	//DESC - Db Connection
	db := database.Conn.DB
	//DESC - Emplouee Model
	seller := models.Seller{}
	//DESC - user model for check mail
	user := models.User{Email: loginInput.Email}
	//DESC - Check validation for email
	validateError := []map[string]interface{}{}
	validateError = globals.ValidateStruct(user)
	if len(validateError) > 0 {
		return nil, errors.New("Email is not valid")
	}

	//DESC - Check user is exist
	err = db.Where(&user).First(&user).Error
	if user.ID == 0 || err != nil {
		return nil, errors.New("User not found or password is wrong")
	}
	// if !globals.CheckPasswordHash(user.Password, []byte(loginInput.Password)) {
	// 	fmt.Printf("Invalid password!!!")
	// 	return nil, errors.New("Invalid password")
	// }

	//DESC - Get employee identification,ability,userLevel and roles by user id
	seller.UserID = user.ID
	err = db.Preload("User").Preload("User.UserRole").Preload("Identification").Where(&seller).First(&seller).Error
	if err != nil || seller.ID == 0 || seller.UserID == 0 {
		return nil, errors.New("No such employee was found.")
	}

	GetAbilityByUserID(&seller.User)

	return &seller, nil
}

func (Authentication) Login(c *fiber.Ctx) error {
	var loginInput LoginInput

	if err := c.BodyParser(&loginInput); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(globals.Response{
			Error:   true,
			Message: "Invalid request",
		})
	}
	//DESC - Bringing your information after the user control
	seller, err := getUserInfo(&loginInput)

	if err != nil {
		c.Status(401)
		return c.JSON(globals.Response{
			Error:   true,
			Message: err.Error(),
		})
	}
	//DESC - GET token time in ENV file
	tokenTime := time.Now().Add(time.Duration(secret.Env["JWT"].(map[string]any)["expiresIn"].(float64)) * time.Minute)
	//DESC - Create token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"seller_name": seller.Identification.Name + " " + seller.Identification.Surname,
		"user_id":     seller.User.ID,
		"seller_id":   seller.ID,
		"level":       seller.User.Level,
		"product":     "auction",
		"ExpiresAt":   tokenTime.Unix(), //Data in Env
	})

	var token string
	//DESC - Token secret key
	token, err = claims.SignedString([]byte(secret.Env["JWT"].(map[string]any)["secret"].(string)))
	if err != nil {
		c.Status(400).JSON(globals.Response{
			Error:   true,
			Message: "Error while creating token",
		})
	}
	sellerData := UserData{
		SellerName: seller.Identification.Name + " " + seller.Identification.Surname,
		UserID:     seller.User.ID,
		SellerID:   seller.ID,
		Level:      seller.User.Level,
		ExpiresAt:  tokenTime.Unix(),
		Token:      token,
	}

	return c.Status(fiber.StatusOK).JSON(globals.Response{
		Body:    sellerData,
		Message: "Login success",
	})
}

func (Authentication) Register(c *fiber.Ctx) error {
	//seller register
	var registerInput RegisterInput
	if err := c.BodyParser(&registerInput); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(globals.Response{
			Error:   true,
			Message: "Invalid request",
		})
	}
	//DESC - Check validation for email
	validateError := []map[string]interface{}{}
	validateError = globals.ValidateStruct(registerInput)
	if len(validateError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(globals.Response{
			Error:   true,
			Body:    validateError,
			Message: "Invalid request",
		})
	}

	//DESC - Check user is exist
	user := models.User{Email: registerInput.Email}
	err := database.Conn.DB.Where(&user).First(&user).Error
	if user.ID != 0 || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(globals.Response{
			Error:   true,
			Message: "User already exist",
		})
	}

	pass, err := models.HashPassword([]byte(registerInput.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(globals.Response{
			Error:   true,
			Message: "Password could not be created",
		})
	}

	//DESC - Create user
	user = models.User{
		Email:    registerInput.Email,
		Password: pass,
		Level:    1,
	}
	err = database.Conn.DB.Create(&user).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(globals.Response{
			Error:   true,
			Message: "User could not be created",
		})
	}

	//DESC - Create identification
	identification := models.Identification{
		SellerID: user.ID,
		Name:     registerInput.Name,
		Surname:  registerInput.Surname,
	}
	err = database.Conn.DB.Create(&identification).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(globals.Response{
			Error:   true,
			Message: "Identification could not be created",
		})
	}
	//DESC - Create seller
	seller := models.Seller{
		UserID:         user.ID,
		Identification: identification,
	}
	err = database.Conn.DB.Create(&seller).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(globals.Response{
			Error:   true,
			Message: "Seller could not be created",
		})
	}

	return c.Status(fiber.StatusOK).JSON(globals.Response{
		Message: "Register success",
	})

}
