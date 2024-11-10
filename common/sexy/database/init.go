package database

func GetAllDBTable() (dst []interface{}) {
	dst = append(dst, &Account{})
	dst = append(dst, &AccountActive{})
	dst = append(dst, &AccountCollect{})
	dst = append(dst, &Project{})
	dst = append(dst, &ProjectLike{})
	dst = append(dst, &ProjectSuperLike{})
	dst = append(dst, &ProjectUnLike{})
	return
}
