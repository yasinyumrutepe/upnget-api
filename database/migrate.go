package database

import "auction/models"

type MigrateParams struct {
	Model     interface{}
	Field     string
	JoinTable interface{}
}

var (
	migrateRelationList = []MigrateParams{
		{&models.User{}, "Abilities", &models.UserAbilities{}},
	}

	// *Financel_info,
	// {model.FinancelInfo{}, "FinancelInfo", &model.Personal{}},
	//{&m.Financel_info{}, "financel_info", &m.Personal{}},

	// 	// 	// //* Market
	// 	// 	{&m.Market{}, "Markups", &m.MarketMarkups{}},
	// }
	//migrateModelList = []interface{}{models.FinancelInfo{}, models.Personal{}} daha öcneki many2many için

	migrateModelList = []interface{}{
		&models.Seller{},
		&models.Identification{},
		&models.Module{},
		&models.User{},
		&models.Address{},
		&models.Brand{},
		&models.Contract{},
		&models.OrderType{},
		&models.Order{},
		&models.Bid{},
		&models.Category{},
		&models.Product{},
		&models.ShippingCompany{},
		&models.PaymentApi{},
		&models.Advertising{},
		&models.File{},
	}

	// seederModelList = []globals.Seeder{
	// 	// &m.AirportCodes{},
	// 	// &m.Term{},
	// }
)
