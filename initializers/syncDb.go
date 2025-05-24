package initializers

import (
	"fmt"

	"github.com/Divyshekhar/finsnap-go/models"
)

func SyncDb() {
	Db.AutoMigrate(&models.User{}, &models.Expense{}, &models.Income{})
	fmt.Println("Migrated")
}
