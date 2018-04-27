package scripts

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/sudarshan-reddy/benjerry/models"
)

//MoveData moves the data from icecream.json to the db
func MoveData(iceCreamStore models.IceCreamStore) error {
	f, err := os.Open("icecream.json")
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	re := regexp.MustCompile("[[:^ascii:]]")
	t := re.ReplaceAllLiteralString(string(data), "")

	var iceCreamData []models.IceCream

	err = json.Unmarshal([]byte(t), &iceCreamData)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return iceCreamStore.WithTxContext(ctx, func(ctx context.Context) error {
		for _, iceCreamInfo := range iceCreamData {
			err := iceCreamStore.StoreContext(ctx, iceCreamInfo)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
