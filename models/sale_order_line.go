package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"goERP/utils"

	"github.com/astaxie/beego/orm"
)

// SaleOrderLine 订单明细
type SaleOrderLine struct {
	ID            int64               `orm:"column(id);pk;auto" json:"id"`              //主键
	CreateUser    *User               `orm:"rel(fk);null" json:"-"`                //创建者
	UpdateUser    *User               `orm:"rel(fk);null" json:"-"`                //最后更新者
	CreateDate    time.Time           `orm:"auto_now_add;type(datetime)" json:"-"` //创建时间
	UpdateDate    time.Time           `orm:"auto_now;type(datetime)" json:"-"`     //最后更新时间
	FormAction    string              `orm:"-" form:"FormAction"`                  //非数据库字段，用于表示记录的增加，修改
	Name          string              `orm:"default(\"\")" json:"name"`            //订单明细号
	Company       *Company            `orm:"rel(fk)"`                              //公司
	SaleOrder     *SaleOrder          `orm:"rel(fk);null" `                        //销售订单
	Partner       *Partner            `orm:"rel(fk)"`                              //客户
	Product       *ProductProduct     `orm:"rel(fk)"`                              //产品
	FirstSaleUom  *ProductUom         `orm:"rel(fk)"`                              //第一销售单位
	SecondSaleUom *ProductUom         `orm:"rel(fk)"`                              //第二销售单位
	FirstSaleQty  float32             `orm:"default(1)"`                           //第一销售单位
	SecondSaleQty float32             `orm:"default(0)"`                           //第二销售单位
	State         *SaleOrderLineState `orm:"rel(fk)"`                              //订单明细状态
}

func init() {
	orm.RegisterModel(new(SaleOrderLine))
}

// AddSaleOrderLine insert a new SaleOrderLine into database and returns
// last inserted ID on success.
func AddSaleOrderLine(obj *SaleOrderLine) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(obj)
	return id, err
}

// GetSaleOrderLineByID retrieves SaleOrderLine by ID. Returns error if
// ID doesn't exist
func GetSaleOrderLineByID(id int64) (obj *SaleOrderLine, err error) {
	o := orm.NewOrm()
	obj = &SaleOrderLine{ID: id}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetAllSaleOrderLine retrieves all SaleOrderLine matches certain condition. Returns empty list if
// no records exist
func GetAllSaleOrderLine(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (utils.Paginator, []SaleOrderLine, error) {
	var (
		objArrs   []SaleOrderLine
		paginator utils.Paginator
		num       int64
		err       error
	)
	o := orm.NewOrm()
	qs := o.QueryTable(new(SaleOrderLine))
	qs = qs.RelatedSel()

	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return paginator, nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return paginator, nil, errors.New("Error: unused 'order' fields")
		}
	}

	qs = qs.OrderBy(sortFields...)
	if cnt, err := qs.Count(); err == nil {
		paginator = utils.GenPaginator(limit, offset, cnt)
	}
	if num, err = qs.Limit(limit, offset).All(&objArrs, fields...); err == nil {
		paginator.CurrentPageSize = num
	}
	return paginator, objArrs, err
}

// UpdateSaleOrderLineByID updates SaleOrderLine by ID and returns error if
// the record to be updated doesn't exist
func UpdateSaleOrderLineByID(m *SaleOrderLine) (err error) {
	o := orm.NewOrm()
	v := SaleOrderLine{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// GetSaleOrderLineByName retrieves SaleOrderLine by Name. Returns error if
// Name doesn't exist
func GetSaleOrderLineByName(name string) (obj *SaleOrderLine, err error) {
	o := orm.NewOrm()
	obj = &SaleOrderLine{Name: name}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// DeleteSaleOrderLine deletes SaleOrderLine by ID and returns error if
// the record to be deleted doesn't exist
func DeleteSaleOrderLine(id int64) (err error) {
	o := orm.NewOrm()
	v := SaleOrderLine{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SaleOrderLine{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}