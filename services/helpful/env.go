package helpful

import (
	"log"
	"os"
)

func GetEnvParam(param string, mandatory bool) string {
	value := os.Getenv(param)
	if mandatory {
		if value == "" {
			log.Fatalf("[ERROR] The env parameter %s is not found", param)
		}
	}

	return value
}
