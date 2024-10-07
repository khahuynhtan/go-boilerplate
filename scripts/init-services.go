package scripts

import (
	auth "myapp/api/author/services"
)

func InitServices() {
	auth.InitService()
}
