package model

// Category 属性后面的 单引号括起来的字符串，叫 tag，其实相当于 java 里的注解
// json tag 用来表示 序列化成 json 后，key 对应的名称
// gorm 表示 该属性对应的建表语句
type Category struct {
	/*
		这个属性，相当于 java 里的

			@json("id")
			@gorm("primary_key")
			int id;
	*/
	ID uint `json:"id" gorm:"primary_key"`

	Name string `json:"name" gorm:"type:varchar(50);not null;unique"`

	CreatedAt Time `json:"created_at" gorm:"type:timestamp default CURRENT_TIMESTAMP"`

	UpdatedAt Time `json:"updated_at" gorm:"type:timestamp default CURRENT_TIMESTAMP"`
}

/*

class Category {

	@json("id")
	@gorm("primary_key")
	public int ID;

	@json("name")
	@gorm("type:varchar(50);not null;unique")
	public int Name;

	@json("created_at")
	@gorm("type:timestamp default CURRENT_TIMESTAMP")
	public int CreatedAt;

	@json("updated_at")
	@gorm("type:timestamp default CURRENT_TIMESTAMP")
	public int UpdatedAt;

}
*/
