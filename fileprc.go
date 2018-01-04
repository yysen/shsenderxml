package main

import (
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

type Result struct {
	Error   bool
	Message string
}
type MyTime struct {
	time.Time
}

func (c MyTime) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return []byte("null"), nil
	}
	str, err := json.Marshal(c.Time)
	return str, err
}
func (c *MyTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02T15:04:05" // yyyymmdd date format
	var v string
	d.DecodeElement(&v, &start)
	if len(v) == 0 {
		return nil
	}
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = MyTime{parse}
	return nil
}

type MyFloat float64

func (c *MyFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	if len(v) == 0 {
		return nil
	}
	parse, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	*c = MyFloat(parse)
	return nil
}

type Corporation_info_investor struct {
	Investor_nationality string  `xml:"investor_nationality"` //股东国别代码
	Subconprop           MyFloat `xml:"subconprop"`           //认缴出资比例
	Blicno               string  `xml:"blicno"`               //证照号码
	Conform              string  `xml:"conform"`              //认缴出资方式
	Cer_type             string  `xml:"cer_type"`             //证件类型
	Subconam             MyFloat `xml:"subconam"`             //认缴出资额
	Condate              MyTime  `xml:"condate"`              //认缴出资时间
	Cer_no               string  `xml:"cer_no"`
	Ent_name             string  `xml:"ent_name"` //名称
	Investor             string  `xml:"investor"` //股东姓名
	Actconam             MyFloat `xml:"actconam"` //实缴出资额
	Blictype             string  `xml:"blictype"` //证照类型
}

type Corporation_finance_person struct {
	Finance_name string `xml:"finance_name"` //负责人姓名
	Tel          string `xml:"tel"`          //固定电话
	Email        string `xml:"email"`        //电子邮箱
	Cer_type     string `xml:"cer_type"`     //证件类型
	Cer_no       string `xml:"cer_no"`       //证件号码
	Mob_tel      string `xml:"mob_tel"`      //移动电话
}
type Corporation_info_change struct {
	Change_item   string `xml:"change_item"`
	Change_benote string `xml:"change_benote"`
	Change_afnote string `xml:"change_afnote"`
	Change_date   MyTime `xml:"change_date"`
}
type Corporation_info_reg struct {
	Person_name          string  `xml:"person_name"`          //法定代表人
	Telephone            string  `xml:"telephone"`            //联系电话
	Corp_type            string  `xml:"corp_type"`            //法人类型
	Establish_date       MyTime  `xml:"establish_date"`       //开业（设立）日期
	Currency             string  `xml:"currency"`             //币种
	Receiving_organ      string  `xml:"receiving_organ"`      //受理机关代码
	Area_code            string  `xml:"area_code"`            //所属行政区划代码
	Uni_sc_id            string  `xml:"uni_sc_id"`            //统一社会信用代码
	Organ_code           string  `xml:"organ_code"`           //组织机构代码
	Zip                  string  `xml:"zip"`                  //邮政编码
	Person_cert_code     string  `xml:"person_cert_code"`     //个人身份证号码
	Person_landline_tel  string  `xml:"person_landline_tel"`  //个人固定电话
	Corp_info_id         string  `xml:"corp_info_id"`         //数据中心法人实体序号
	Reg_capital          MyFloat `xml:"reg_capital"`          //注册资本/注册资金/资金数额/出资额/成员出资总额
	Corp_status          string  `xml:"corp_status"`          //法人状态
	Reg_upd_date         MyTime  `xml:"reg_upd_date"`         //登记类业务发布时间
	Sender_nationality   string  `xml:"sender_nationality"`   //国家（地区）
	Person_email         string  `xml:"person_email"`         //个人邮箱
	Corp_name            string  `xml:"corp_name"`            //名称
	Entity_id            string  `xml:"entity_id"`            //工商标识
	Sender_name          string  `xml:"sender_name"`          //外国（地区）企业名称
	Industry_code        string  `xml:"industry_code"`        //行业代码
	Person_cert_type     string  `xml:"person_cert_type"`     //身份证件类型
	Address              string  `xml:"address"`              //地址
	Business_scope       string  `xml:"business_scope"`       //经营范围
	Reg_no               string  `xml:"reg_no"`               //注册号
	Change_date          MyTime  `xml:"change_date"`          //变更日期
	Sender_address       string  `xml:"sender_address"`       //外国（地区）企业住所
	Contact_name         string  `xml:"contact_name"`         //联络员姓名
	Etps_type            string  `xml:"etps_type"`            //公司类型/经济性质/企业类型（登记注册类型）
	Parent_uni_sc_id     string  `xml:"parent_uni_sc_id"`     //分支机构上级主管部门（公司、隶属单位、企业）统一社会信用代码
	Pro_loc              string  `xml:"pro_loc"`              //生产经营地址
	Revoke_date          MyTime  `xml:"revoke_date"`          //吊销日期
	Parent_name          string  `xml:"parent_name"`          //分支机构上级主管部门（公司、隶属单位、企业）名称
	Calc_method          string  `xml:"calc_method"`          //核算方式
	Parent_location_type string  `xml:"parent_location_type"` //上级法人详细地址
	Reg_organ            string  `xml:"reg_organ"`            //企业登记机关
	Trd_end_date         MyTime  `xml:"trd_end_date"`         //经营期限/营业期限/合伙期限自
	Trd_start_date       MyTime  `xml:"trd_start_date"`       //经营期限/营业期限/合伙期限止
	Parent_reg_no        string  `xml:"parent_reg_no"`        //分支机构上级主管部门（公司、隶属单位、企业）注册号
	Person_telephone     string  `xml:"person_telephone"`     //法定代表人（负责人）固定电话
	Total_investment     MyFloat `xml:"total_investment"`     //投资总额
}
type entInfo struct {
	retype                 string
	GUID_01                string
	UNISCID_01             string //统一社会信用代码
	BUSITYPE_01            string //工商业务类型
	ORGNCODE_01            string //组织机构代码
	REGNO_01               string //注册号
	ENTNAME_01             string //名称
	ENTTYPE_01             string //类型
	UNITTYPE_01            string //单位类型
	REGTYPE_01             string //登记注册类型
	HOLDINGSMSG_01         string //控股情况
	ORGTYPE_01             string //机构类型
	ACCCATE_01             string //执行会计标准类别
	REGORG_01              string //企业登记机关
	AREACODE_01            string //数据处理地代码
	ESTDATE_01             string //开业（设立）日期，格式yyyyMMdd
	OPFROM_01              string //经营期限/营业期限/合伙期限自，格式yyyyMMdd
	OPTO_01                string //经营期限/营业期限/合伙期限止，格式yyyyMMdd
	OPERATIONPERIOD_01     string //经营期限
	DOM_01                 string //住所
	POSTCODE_01            string //邮政编码
	TEL_01                 string //联系电话
	PROLOC_01              string //生产经营地址
	OPSCOPE_01             string //经营范围
	REGCAP_01              string //注册资本
	REGCAPCUR_01           string //货币种类
	CURAMOUNT_01           string //币种金额
	INFOACTIONTYPE_01      string //信息操作类型
	S_MODIFYTIME_01        string //数据修改时间，格式为yyyyMMddHHmmss
	S_UNITCODE_01          string //行政区划代码
	S_SJBBH_01             string //数据包编码
	S_ISCANCEL_01          string //是否被注销注销
	S_BATCHGUID_01         string //批次号，数据上传所属批次号，用于判断是否是同一批次数据
	S_PARENTENTDOCUMENT_01 string //上级主管部门名称
	S_LEREP_01             string //法定代表人姓名
	S_INFOFINANCE_01       string //财务负责人姓名
	S_INVESTORCOUNT_01     string //投资人数量
	S_CHILDENTCOUNT_01     string //下级子公司数量
	S_ALTDATE_01           string //变更信息最新一次变更的时间，取自变更信息中变更时间
	S_STATUS_01            string //状态，999：正常，100：必要性审核，200：核实性审核
	S_ISREPEAT_01          string //统一社会信用代码是否重复，0：不重码，1：重码
	S_HANDLETYPE_01        string //对重码数据的人工处理，0：未处理，1：正常，2：删除
	S_ISAUDITED_01         string //是否已审核，0：未审核，1：已审核
	S_ISPUSHEDMLK_01       string //是否已推送到名录库，0：未推送，1：推送
	S_UPLOADTIME_01        string //数据上传时间，格式为yyyyMMddHHmmss
	S_AUDITTIME_01         string //数据审核的时间，格式为yyyyMMddHHmmss
}

//法人信息表
type faren struct {
	retype     string
	GUID_01    string
	GUID_02    string //唯一主键guid
	LEREP_02   string //法人姓名
	CERTYPE_02 string //证件类型
	CERNO_02   string //证件编号
	TEL_02     string //固话
	MOBTEL_02  string //移动电话
	EMAIL_02   string //电子邮件

}

//财务负责人
type caiwu struct {
	retype     string
	GUID_01    string
	GUID_03    string //唯一主键guid
	NAME_03    string //负责人姓名
	CERTYPE_03 string //证件类型
	CERNO_03   string //证件编号
	TEL_03     string //固话
	MOBTEL_03  string //移动电话
	EMAIL_03   string //电子邮件
}

//上级主管部门信息表
type shangji struct {
	retype        string
	GUID_01       string
	GUID_04       string //唯一主键guid
	UNISCID_04    string //主管部门统一社会信用代码
	BRORGNCODE_04 string //上级组织机构代码
	BRNAME_04     string //主管部门名称
	REGNO_04      string //注册号
}

//投资人信息
type touzhi struct {
	retype        string
	GUID_01       string
	GUID_05       string //唯一主键guid
	INVTYPE_05    string //投资类型
	INV_05        string //投资方名称
	CERTYPE_05    string //证件类型
	CERNO_05      string //证件编号
	BLICTYPE_05   string //证照类型
	BLICNO_05     string //证照号码
	SUBCONAM_05   string //认缴出资额
	CONDATE_05    string //认缴出资时间，格式yyyyMMdd
	CONFORM_05    string //认缴出资方式
	SUBCONPROP_05 string //认缴出资比例
	COUNTRY_05    string //国籍
	CURRENCY_05   string //币种
}

//下级子公司信息表
type xiaji struct {
	retype      string
	GUID_01     string
	GUID_06     string //唯一主键guid
	UNISCID_06  string //分支机构统一社会信用代码
	REGNO_06    string //注册号
	BRNAME_06   string //分支机构名称
	ENTTYPE_06  string //类型
	ESTDATE_06  string //成立日期，格式yyyyMMdd
	PRIL_06     string //负责人
	OPLOC_06    string //营业场所
	OPSCOPE_06  string //营业期限自，格式yyyyMMdd
	OPTO_06     string //营业期限至，格式yyyyMMdd
	REGORG_06   string //登记机关
	APPRDATE_06 string //核准日期，格式yyyyMMdd
}

//变更信息表
type biangen struct {
	retype     string
	GUID_01    string
	GUID_07    string //唯一主键guid
	ALTITEM_07 string //变更事项
	ALTBE_07   string //变更前内容
	ALTAF_07   string //变更后内容
	ALTDATE_07 string //变更日期，格式yyyyMMdd
	S_SJBBH_07 string //数据包编号，此条变更信息上传时的数据包编号
}

//注销信息
type zhuxiao struct {
	retype          string
	GUID_01         string
	GUID_08         string //唯一主键guid
	CANDATE_08      string //注销日期，格式yyyyMMdd
	EQUPLECANREA_08 string //注销原因
}

//表示多张表的一条记录
type Corp_info struct {
	Corporation_info_reg       Corporation_info_reg        `xml:"corporation_info_reg"`
	Corporation_finance_person Corporation_finance_person  `xml:"corporation_finance_person"`
	Corporation_info_investor  []Corporation_info_investor `xml:"corporation_info_investor"`
	Corporation_info_change    []Corporation_info_change   `xml:"corporation_info_change"`
}

//表示一个xml文件
type Project struct {
	Corp_info []Corp_info `xml:"corp_info"`
}

var (
	entSize = []int{2, 32, 18, 1, 9, 50, 100, 4, 1, 3, 1, 2, 1, 6, 15, 20, 20, 20, 6, 200, 6, 60, 200, 3000, 24, 3, 24, 1, 14, 40, 20, 1, 32, 100, 100, 50, 6, 6, 20, 3, 1, 1, 1, 1, 14, 14}
	frSize  = []int{2, 32, 32, 100, 2, 60, 60, 60, 100}
	cwSize  = []int{2, 32, 32, 50, 2, 60, 30, 30, 100}
	sjSize  = []int{2, 32, 32, 18, 9, 100, 50}
	tzSize  = []int{2, 32, 32, 10, 200, 2, 60, 10, 60, 14, 20, 100, 10, 30, 10}
	xjSize  = []int{2, 32, 32, 18, 50, 100, 4, 20, 30, 200, 3000, 20, 20, 6, 20}
	bgSize  = []int{2, 32, 32, 100, 4000, 4000, 20, 20}
	zxSize  = []int{2, 32, 32, 20, 200}
)

//遍历结构体，并处理相关字段
func toStringStruct(stru interface{}, fieldSize []int) []string {
	strs := []string{}
	//得到stru的结构
	value := reflect.ValueOf(stru)
	//遍历结构体
	for i := 0; i < value.NumField(); i++ {
		st := value.Field(i).String()
		strs = append(strs, changeLen(st, fieldSize[i]))
	}
	return strs
}

//将数据处理成国家需要的长度,是字节的长度
func changeLen(data string, size int) string {
	//去除换行符
	str := strings.Replace(strings.Replace(data, "\r", " ", -1), "\n", " ", -1)
	//得到字符数组
	rstr := []rune(str)
	//得到字符的长度，而不是字节，防止出现乱码
	dataLen := len(rstr)
	if dataLen <= size {
		str = str + strings.Repeat(" ", size-dataLen)
	} else if dataLen > size {
		str = string(rstr[:size])
	}
	return str
}

//读取一条主表及其子表记录，并将结果转化为国家需要的格式
func readOneFile(fs string) ([][]string, error) {
	//读取fs文件里的内容
	filebys, err := ioutil.ReadFile(fs)
	if err != nil {
		logOut.Error("read file error")
		return nil, err
	}
	prj := &Project{}
	if err = xml.Unmarshal(filebys, prj); err != nil {
		logOut.Error("Unmarshal error", err)
		return nil, err
	}
	if len(prj.Corp_info) != 1 {
		logOut.Error("corp_info not is one record")
	}
	//读取的数据
	datas := [][]string{}
	//主表uuid
	entUUID := hex.EncodeToString(uuid.NewUUID())
	ENTTYPE, REGTYPE, HOLDINGSMSG, UNITTYPE, ORGTYPE, ACCCATE := transform(etpsTypeDict, prj.Corp_info[0].Corporation_info_reg.Etps_type)
	entInfo := entInfo{
		retype:  "01",
		GUID_01: entUUID,
		UNISCID_01: func() string {
			scid := prj.Corp_info[0].Corporation_info_reg.Uni_sc_id
			if len(scid) == 0 {
				return "--"
			}
			return scid
		}(), //社会信用代码
		BUSITYPE_01: func() string {
			if len(prj.Corp_info[0].Corporation_info_change) > 0 {
				return "2"
			} else if !prj.Corp_info[0].Corporation_info_reg.Revoke_date.IsZero() {
				return "6"
			}
			return "1"
		}(), //工商业务类型
		ORGNCODE_01: prj.Corp_info[0].Corporation_info_reg.Organ_code, //组织机构代码
		REGNO_01: func() string {
			scid := prj.Corp_info[0].Corporation_info_reg.Reg_no
			if len(scid) == 0 {
				return "--"
			}
			return scid
		}(), //注册号
		ENTNAME_01:     prj.Corp_info[0].Corporation_info_reg.Corp_name, //名称
		ENTTYPE_01:     ENTTYPE,                                         //类型
		UNITTYPE_01:    UNITTYPE,                                        //单位类型
		REGTYPE_01:     REGTYPE,                                         //登记注册类型
		HOLDINGSMSG_01: HOLDINGSMSG,                                     //控股情况
		ORGTYPE_01:     ORGTYPE,                                         //机构类型
		ACCCATE_01:     ACCCATE,                                         //执行会计标准类别
		REGORG_01: func() string {
			reg := strings.Replace(prj.Corp_info[0].Corporation_info_reg.Reg_organ, " ", "", -1)
			regs := convert(areacodeDict, reg)
			if len(regs) == 6 {
				return regs
			} else {
				//判断是否有汉字
				rege := regexp.MustCompile("[\\p{Han}]+")
				if len(rege.FindAllString(reg, -1)) > 0 {
					return ""
				}
			}
			return reg
		}(), //企业登记机关
		AREACODE_01: convert(areacodeDict, prj.Corp_info[0].Corporation_info_reg.Area_code), //数据处理地代码
		ESTDATE_01: func() string {
			dt := prj.Corp_info[0].Corporation_info_reg.Establish_date
			if dt.IsZero() {
				return ""
			}
			return dt.Format("20060102")

		}(), //开业（设立）日期，格式yyyyMMdd
		OPFROM_01: func() string {
			dt := prj.Corp_info[0].Corporation_info_reg.Trd_start_date
			if dt.IsZero() {
				return ""
			}
			return dt.Format("20060102")

		}(), //经营期限/营业期限/合伙期限自，格式yyyyMMdd
		OPTO_01: func() string {
			dt := prj.Corp_info[0].Corporation_info_reg.Trd_end_date
			if dt.IsZero() {
				return ""
			}
			return dt.Format("20060102")

		}(), //经营期限/营业期限/合伙期限止，格式yyyyMMdd
		OPERATIONPERIOD_01: "",                                                                                          //经营期限
		DOM_01:             prj.Corp_info[0].Corporation_info_reg.Address,                                               //住所
		POSTCODE_01:        prj.Corp_info[0].Corporation_info_reg.Zip,                                                   //邮政编码
		TEL_01:             prj.Corp_info[0].Corporation_info_reg.Telephone,                                             //联系电话
		PROLOC_01:          prj.Corp_info[0].Corporation_info_reg.Pro_loc,                                               //生产经营地址
		OPSCOPE_01:         prj.Corp_info[0].Corporation_info_reg.Business_scope,                                        //经营范围
		REGCAP_01:          strconv.FormatFloat(float64(prj.Corp_info[0].Corporation_info_reg.Reg_capital), 'f', 6, 64), //注册资本
		REGCAPCUR_01:       convert(curDict, prj.Corp_info[0].Corporation_info_reg.Currency),                            //货币种类
		CURAMOUNT_01:       strconv.FormatFloat(float64(prj.Corp_info[0].Corporation_info_reg.Reg_capital), 'f', 6, 64), //币种金额
		INFOACTIONTYPE_01: func() string {
			if len(prj.Corp_info[0].Corporation_info_change) > 0 {
				return "1"
			} else if !prj.Corp_info[0].Corporation_info_reg.Revoke_date.IsZero() {
				return "2"
			}
			return "0"
		}(), //信息操作类型
		S_MODIFYTIME_01:        "",                                                                     //数据修改时间，格式为yyyyMMddHHmmss
		S_UNITCODE_01:          convert(areacodeDict, prj.Corp_info[0].Corporation_info_reg.Area_code), //行政区划代码
		S_SJBBH_01:             "",                                                                     //数据包编码
		S_ISCANCEL_01:          "",                                                                     //是否被注销注销
		S_BATCHGUID_01:         "",                                                                     //批次号，数据上传所属批次号，用于判断是否是同一批次数据
		S_PARENTENTDOCUMENT_01: prj.Corp_info[0].Corporation_info_reg.Parent_name,                      //上级主管部门名称
		S_LEREP_01:             prj.Corp_info[0].Corporation_info_reg.Person_name,                      //法定代表人姓名
		S_INFOFINANCE_01:       prj.Corp_info[0].Corporation_finance_person.Finance_name,               //财务负责人姓名
		S_INVESTORCOUNT_01:     fmt.Sprintf("%d", len(prj.Corp_info[0].Corporation_info_change)),       //投资人数量
		S_CHILDENTCOUNT_01:     "",                                                                     //下级子公司数量
		S_ALTDATE_01:           "",                                                                     //变更信息最新一次变更的时间，取自变更信息中变更时间
		S_STATUS_01:            "",                                                                     //状态，999：正常，100：必要性审核，200：核实性审核
		S_ISREPEAT_01:          "",                                                                     //统一社会信用代码是否重复，0：不重码，1：重码
		S_HANDLETYPE_01:        "",                                                                     //对重码数据的人工处理，0：未处理，1：正常，2：删除
		S_ISAUDITED_01:         "",                                                                     //是否已审核，0：未审核，1：已审核
		S_ISPUSHEDMLK_01:       "",                                                                     //是否已推送到名录库，0：未推送，1：推送
		S_UPLOADTIME_01:        "",                                                                     //数据上传时间，格式为yyyyMMddHHmmss
		S_AUDITTIME_01:         "",                                                                     //数据审核的时间，格式为yyyyMMddHHmmss
	}
	datas = append(datas, toStringStruct(entInfo, entSize))
	//法人信息表
	fr := faren{
		retype:     "02",
		GUID_01:    entUUID,
		GUID_02:    hex.EncodeToString(uuid.NewUUID()),                                            //唯一主键guid
		LEREP_02:   prj.Corp_info[0].Corporation_info_reg.Person_name,                             //法人姓名
		CERTYPE_02: convert(certTypeDict, prj.Corp_info[0].Corporation_info_reg.Person_cert_type), //证件类型
		CERNO_02:   prj.Corp_info[0].Corporation_info_reg.Person_cert_code,                        //证件编号
		TEL_02: func() string {
			tel := prj.Corp_info[0].Corporation_info_reg.Person_telephone
			match, _ := regexp.MatchString("^0[0-9]", tel)
			if match && (len(tel) < 11) {
				return ""
			}
			return tel
		}(), //固话
		MOBTEL_02: func() string {
			tel := prj.Corp_info[0].Corporation_info_reg.Person_telephone
			match, _ := regexp.MatchString("^1[0-9]{10}", tel)
			if match {
				return tel
			}
			return ""
		}(), //移动电话
		EMAIL_02: prj.Corp_info[0].Corporation_info_reg.Person_email, //电子邮件

	}
	//判断有无法人信息
	if len(fr.CERNO_02+fr.CERTYPE_02+fr.EMAIL_02+fr.LEREP_02+fr.TEL_02) > 0 {
		datas = append(datas, toStringStruct(fr, frSize))
	}
	//财务负责人
	cw := caiwu{
		retype:     "03",
		GUID_01:    entUUID,
		GUID_03:    hex.EncodeToString(uuid.NewUUID()),                                          //唯一主键guiD
		NAME_03:    prj.Corp_info[0].Corporation_finance_person.Finance_name,                    //负责人姓名
		CERTYPE_03: convert(certTypeDict, prj.Corp_info[0].Corporation_finance_person.Cer_type), //证件类型
		CERNO_03:   prj.Corp_info[0].Corporation_finance_person.Cer_no,                          //证件编号
		TEL_03:     "",                                                                          //固话
		MOBTEL_03:  prj.Corp_info[0].Corporation_finance_person.Mob_tel,                         //移动电话
		EMAIL_03:   "",                                                                          //电子邮件
	}
	//判断有无财务信息
	if len(cw.NAME_03+cw.CERTYPE_03+cw.CERNO_03+cw.MOBTEL_03) > 0 {
		datas = append(datas, toStringStruct(cw, cwSize))
	}
	//上级主管部门信息表
	sj := shangji{
		retype:     "04",
		GUID_01:    entUUID,
		GUID_04:    hex.EncodeToString(uuid.NewUUID()),                     //唯一主键guid
		UNISCID_04: prj.Corp_info[0].Corporation_info_reg.Parent_uni_sc_id, //主管部门统一社会信用代码
		BRORGNCODE_04: func() string {
			if len(prj.Corp_info[0].Corporation_info_reg.Parent_uni_sc_id) < 8 {
				return ""
			}
			br := []rune(prj.Corp_info[0].Corporation_info_reg.Parent_uni_sc_id)
			return string(br[8:17])
		}(), //上级组织机构代码
		BRNAME_04: prj.Corp_info[0].Corporation_info_reg.Parent_name, //主管部门名称
		REGNO_04:  prj.Corp_info[0].Corporation_info_reg.Reg_no,      //注册号
	}
	//判断有无上级主管部门信息
	if len(sj.UNISCID_04+sj.BRNAME_04+sj.REGNO_04) > 0 {
		datas = append(datas, toStringStruct(sj, sjSize))
	}
	investor := prj.Corp_info[0].Corporation_info_investor
	//投资人信息
	for _, oneInvestor := range investor {
		tz := touzhi{
			retype:      "05",
			GUID_01:     entUUID,
			GUID_05:     hex.EncodeToString(uuid.NewUUID()),                             //唯一主键guid
			INVTYPE_05:  convert(tzTypeDict, oneInvestor.Blictype),                      //投资类型
			INV_05:      oneInvestor.Investor,                                           //投资方名称
			CERTYPE_05:  convert(certTypeDict, oneInvestor.Cer_type),                    //证件类型
			CERNO_05:    oneInvestor.Cer_no,                                             //证件编号
			BLICTYPE_05: convert(blicTypeDict, oneInvestor.Blictype),                    //证照类型
			BLICNO_05:   oneInvestor.Blicno,                                             //证照号码
			SUBCONAM_05: strconv.FormatFloat(float64(oneInvestor.Subconam), 'f', 6, 64), //认缴出资额
			CONDATE_05: func() string {
				dt := oneInvestor.Condate
				if dt.IsZero() {
					return ""
				}
				return dt.Format("20060102")
			}(), //认缴出资时间，格式yyyyMMdd
			CONFORM_05:    convert(conFormDict, oneInvestor.Conform),                        //认缴出资方式
			SUBCONPROP_05: strconv.FormatFloat(float64(oneInvestor.Subconprop), 'f', 6, 64), //认缴出资比例
			COUNTRY_05:    convert(countryDict, oneInvestor.Investor_nationality),           //国籍
			CURRENCY_05:   convert(curDict, prj.Corp_info[0].Corporation_info_reg.Currency), //币种
		}
		datas = append(datas, toStringStruct(tz, tzSize))
	}

	change := prj.Corp_info[0].Corporation_info_change
	for _, oneChange := range change {
		//变更信息表
		bg := biangen{
			retype:     "07",
			GUID_01:    entUUID,
			GUID_07:    hex.EncodeToString(uuid.NewUUID()), //唯一主键guid
			ALTITEM_07: oneChange.Change_item,              //变更事项
			ALTBE_07:   oneChange.Change_benote,            //变更前内容
			ALTAF_07:   oneChange.Change_afnote,            //变更后内容
			ALTDATE_07: func() string {
				dt := oneChange.Change_date
				if dt.IsZero() {
					return ""
				}
				return dt.Format("20060102")
			}(), //变更日期，格式yyyyMMdd
			S_SJBBH_07: "", //数据包编号，此条变更信息上传时的数据包编号
		}
		datas = append(datas, toStringStruct(bg, bgSize))
	}
	//注销信息
	zx := zhuxiao{
		retype:  "08",
		GUID_01: entUUID,
		GUID_08: hex.EncodeToString(uuid.NewUUID()), //唯一主键guid
		CANDATE_08: func() string {
			dt := prj.Corp_info[0].Corporation_info_reg.Revoke_date
			if dt.IsZero() {
				return ""
			}
			return dt.Format("20060102")
		}(), //注销日期，格式yyyyMMdd
		EQUPLECANREA_08: "", //注销原因
	}
	if len(zx.CANDATE_08) > 0 {
		datas = append(datas, toStringStruct(zx, zxSize))
	}
	return datas, nil
}
