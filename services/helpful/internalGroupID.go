package helpful

import (
	"strconv"
	"strings"
)

func InternalGroupID(groupID int64) int {
	groupIDStr := strconv.Itoa(int(groupID))

	internalGroupID := strings.TrimPrefix(groupIDStr, "-100")

	result, _ := strconv.Atoi(internalGroupID)
	return result
}
