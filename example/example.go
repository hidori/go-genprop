package example

type Example struct {
	value0 int    // ignored
	value1 int    `property:"get"`
	value2 int    `property:"get,set"`
	id     int    `property:"get"`
	api    string `property:"get"`
	url    string `property:"get"`
	http   string `property:"get"`
}
