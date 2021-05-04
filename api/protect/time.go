package protect

import (
	"encoding/json"
	"fmt"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(d []byte) error {
	var u int64
	if err := json.Unmarshal(d, &u); err != nil {
		return err
	}

	*t = Time(time.Unix(u/1000, u%1000*int64(time.Millisecond)))
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	u := time.Time(t).UnixNano()
	return json.Marshal(u / int64(time.Millisecond))
}

func (t Time) GoString() string {
	return fmt.Sprintf("protect.Time(%#v)", time.Time(t).Format(time.RFC3339))
}

type Duration time.Duration

func (d *Duration) UnmarshalJSON(e []byte) error {
	var u int64
	if err := json.Unmarshal(e, &u); err != nil {
		return err
	}

	*d = Duration(u * int64(time.Millisecond))
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(d) / int64(time.Millisecond))
}

func (d Duration) GoString() string {
	return fmt.Sprintf("protect.Duration(%#v)", time.Duration(d).String())
}
