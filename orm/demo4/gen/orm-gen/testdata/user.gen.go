package testdata

import (
    "gitee.com/geektime-geekbang/geektime-go/orm"

    "database/sql"
)

const (
    UserName = "Name"
    UserAge = "Age"
    UserNickName = "NickName"
    UserPicture = "Picture"
)


func UserNameLT(val string) orm.Predicate {
    return orm.C("Name").LT(val)
}

func UserNameGT(val string) orm.Predicate {
    return orm.C("Name").GT(val)
}

func UserNameEQ(val string) orm.Predicate {
    return orm.C("Name").EQ(val)
}

func UserAgeLT(val *int) orm.Predicate {
    return orm.C("Age").LT(val)
}

func UserAgeGT(val *int) orm.Predicate {
    return orm.C("Age").GT(val)
}

func UserAgeEQ(val *int) orm.Predicate {
    return orm.C("Age").EQ(val)
}

func UserNickNameLT(val *sql.NullString) orm.Predicate {
    return orm.C("NickName").LT(val)
}

func UserNickNameGT(val *sql.NullString) orm.Predicate {
    return orm.C("NickName").GT(val)
}

func UserNickNameEQ(val *sql.NullString) orm.Predicate {
    return orm.C("NickName").EQ(val)
}

func UserPictureLT(val []byte) orm.Predicate {
    return orm.C("Picture").LT(val)
}

func UserPictureGT(val []byte) orm.Predicate {
    return orm.C("Picture").GT(val)
}

func UserPictureEQ(val []byte) orm.Predicate {
    return orm.C("Picture").EQ(val)
}


const (
    UserDetailAddress = "Address"
)


func UserDetailAddressLT(val string) orm.Predicate {
    return orm.C("Address").LT(val)
}

func UserDetailAddressGT(val string) orm.Predicate {
    return orm.C("Address").GT(val)
}

func UserDetailAddressEQ(val string) orm.Predicate {
    return orm.C("Address").EQ(val)
}
