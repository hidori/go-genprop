package example

type Example struct {
	value0 int    // ignored
	value1 int    `prop:"get"`
	value2 int    `prop:"get,set"`
	id     int    `prop:"get"`
	api    string `prop:"get"`
	url    string `prop:"get"`
	http   string `prop:"get"`
}
