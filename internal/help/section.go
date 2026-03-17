package help

// Section はヘルプパネルの1グループ分のキーバインド定義。
type Section struct {
	Title string
	Rows  [][2]string // [key label, description]
}
