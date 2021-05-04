package api

import (
	"encoding/json"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(d []byte) error {
	var u int64
	if err := json.Unmarshal(d, &u); err != nil {
		return err
	}

	*t = Time(time.Unix(u, 0))
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	u := time.Time(t).Unix()
	return json.Marshal(u)
}
