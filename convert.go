package main

import (
	"strings"
)

var (
	//公司类型及登记注册类型,控股情况，单位类型，机构类型，执行会计标准
	etpsTypeDict = map[string]string{
		"内资公司":                    "1000    1101",
		"有限责任公司":                  "1100    1101",
		"有限责任公司(国有独资)":            "111015111101",
		"有限责任公司(外商投资企业投资)":        "1120    1101",
		"有限责任公司(外商投资企业合资)":        "112133051101",
		"有限责任公司(外商投资企业与内资合资)":     "1122310 1101",
		"有限责任公司(外商投资企业法人独资)":      "112333051101",
		"有限责任公司(自然人投资或控股)":        "113017331101",
		"有限责任公司(国有控股)":            "114015911101",
		"一人有限责任公司":                "1150    1101",
		"有限责任公司(自然人独资)":           "115117331101",
		"有限责任公司(自然人投资或控股的法人独资)":   "115217331101",
		"有限责任公司(非自然人投资或控股的法人独资)":  "1153    1101",
		"其他有限责任公司":                "1190159 1101",
		"股份有限公司":                  "1200    1101",
		"股份有限公司(上市)":              "1210    1101",
		"股份有限公司(上市、外商投资企业投资)":     "1211340 1101",
		"股份有限公司(上市、自然人投资或控股)":     "121217431101",
		"股份有限公司(上市、国有控股)":         "121316011101",
		"其他股份有限公司(上市)":            "1219160 1101",
		"股份有限公司(非上市)":             "1220    1101",
		"股份有限公司(非上市、外商投资企业投资)":    "1221340 1101",
		"股份有限公司(非上市、自然人投资或控股)":    "122217431101",
		"股份有限公司(非上市、国有控股)":        "122316011101",
		"其他股份有限公司(非上市)":           "1229160 2101",
		"内资分公司":                   "2000    2101",
		"有限责任公司分公司":               "2100    2101",
		"有限责任公司分公司(国有独资)":         "211015112101",
		"有限责任公司分公司(外商投资企业投资)":     "2120    2101",
		"有限责任公司分公司(外商投资企业合资)":     "212133052101",
		"有限责任公司分公司(外商投资企业与内资合资)":  "2122310 2101",
		"有限责任公司分公司(外商投资企业法人独资)":   "212333052101",
		"有限责任公司分公司(自然人投资或控股)":     "213017332101",
		"有限责任公司分公司(国有控股)":         "214015912101",
		"一人有限责任公司分公司":             "2150    2101",
		"有限责任公司分公司(自然人独资)":        "215117332101",
		"有限责任分公司(自然人投资或控股的法人独资)":  "215217332101",
		"有限责任分公司(非自然人投资或控股的法人独资)": "2153    2101",
		"其他有限责任公司分公司":             "2190159 2101",
		"股份有限公司分公司":               "2200    2101",
		"股份有限公司分公司(上市)":           "2210    2101",
		"股份有限公司分公司(上市、外商投资企业投资)":  "2211340 2101",
		"股份有限公司分公司(上市、自然人投资或控股)":  "221217432101",
		"股份有限公司分公司(上市、国有控股)":      "221316012101",
		"其他股份有限公司分公司(上市)":         "2219160 2101",
		"股份有限公司分公司(非上市)":          "2220    2101",
		"股份有限公司分公司(非上市、外商投资企业投资)": "2221340 2101",
		"股份有限公司分公司(非上市、自然人投资或控股)": "222217432101",
		"股份有限公司分公司(国有控股)":         "222316012101",
		"其他股份有限公司分公司(非上市)":        "2229160 2101",
		"内资企业法人":                  "3000    1101",
		"全民所有制":                   "310011011101",
		"集体所有制":                   "320012021101",
		"股份制":                     "3300    1101",
		"股份合作制":                   "3400130 1101",
		"联营":                      "3500    1101",
		"会员制":                     "3700    1101",
		"内资非法人企业、非公司私营企业及内资非公司企业分支机构": "4000        ",
		"事业单位营业":          "4100     20 ",
		"国有事业单位营业":        "4110110  20 ",
		"集体事业单位营业":        "4120120  20 ",
		"社团法人营业":          "4200    140 ",
		"国有社团法人营业":        "4210110 140 ",
		"集体社团法人营业":        "4220120 140 ",
		"内资企业法人分支机构(非法人)": "4300    2101",
		"全民所有制分支机构(非法人)":  "431011012101",
		"集体分支机构(非法人)":     "432012022101",
		"股份制分支机构":         "4330    2101",
		"股份合作制分支机构":       "4340130 2101",
		"经营单位(非法人)":       "4400    2101",
		"国有经营单位(非法人)":     "441011012101",
		"集体经营单位(非法人)":     "442012022101",
		"非公司私营企业":         "4500    1101",
		"合伙企业":            "4530    1101",
		"普通合伙企业":          "453117231101",
		"特殊普通合伙企业":        "453217231101",
		"有限合伙企业":          "453317231101",
		"个人独资企业":          "454017131101",
		"合伙企业分支机构":        "4550    2101",
		"普通合伙企业分支机构":      "455117232101",
		"特殊普通合伙企业分支机构":    "455217232101",
		"有限合伙企业分支机构":      "455317232101",
		"个人独资企业分支机构":      "456017132101",
		//"联营":         "4600 2101",
		"股份制企业(非法人)": "4700    2101",
		"外商投资企业":     "5000    1101",
		//"有限责任公司":              "5100",
		"有限责任公司(中外合资)":        "5110310 1101",
		"有限责任公司(中外合作)":        "5120320 1101",
		"有限责任公司(外商合资)":        "513033051101",
		"有限责任公司(外国自然人独资)":     "514033051101",
		"有限责任公司(外国法人独资)":      "515033051101",
		"有限责任公司(外国非法人经济组织独资)": "516033051101",
		//"股份有限公司":                  "5200",
		"股份有限公司(中外合资、未上市)": "5210340 1101",
		"股份有限公司(中外合资、上市)":  "5220340 1101",
		"股份有限公司(外商合资、未上市)": "523034051101",
		"股份有限公司(外商合资、上市)":  "524034051101",
		"非公司":              "5300    1101",
		"非公司外商投资企业(中外合作)":  "5310320 1101",
		"非公司外商投资企业(外商合资)":  "5320    1101",
		"外商投资合伙企业":         "5400390 1101",
		//"普通合伙企业":                  "5410",
		//"特殊普通合伙企业":                "5420",
		//"有限合伙企业":                  "5430",
		"外商投资企业分支机构":    "5800    2101",
		"分公司":           "581033052101",
		"非公司外商投资企业分支机构": "582033052101",
		"办事处":           "583033052101",
		"外商投资合伙企业分支机构":  "584033052101",
		"台、港、澳投资企业":     "6000    1101",
		//"有限责任公司":                  "6100",
		"有限责任公司(台港澳与境内合资)":     "6110210 1101",
		"有限责任公司(台港澳与境内合作)":     "6120220 1101",
		"有限责任公司(台港澳合资)":        "613023041101",
		"有限责任公司(台港澳自然人独资)":     "614023041101",
		"有限责任公司(台港澳法人独资)":      "615023041101",
		"有限责任公司(台港澳非法人经济组织独资)": "616023041101",
		"有限责任公司(台港澳与外国投资者合资)":  "6170    1101",
		//"股份有限公司":                  "6200",
		"股份有限公司(台港澳与境内合资、未上市)":    "6210240 1101",
		"股份有限公司(台港澳与境内合资、上市)":     "6220240 1101",
		"股份有限公司(台港澳合资、未上市)":       "623024041101",
		"股份有限公司(台港澳合资、上市)":        "624024041101",
		"股份有限公司(台港澳与外国投资者合资、未上市)": "6250    1101",
		"股份有限公司(台港澳与外国投资者合资、上市)":  "6260    1101",
		//"非公司": "6300",
		"非公司台、港、澳企业(台港澳与境内合作)": "6310220 1101",
		"非公司台、港、澳企业(台港澳合资)":    "632023041101",
		"港、澳、台投资合伙企业":          "6400290 1101",
		//"普通合伙企业":               "6410",
		//"特殊普通合伙企业":             "6420",
		//"有限合伙企业":               "6430",
		"台、港、澳投资企业分支机构": "6800    2101",
		//"分公司":                  "6810",
		"非公司台、港、澳投资企业分支机构": "682023042101",
		//"办事处": "6830",
		"港、澳、台投资合伙企业分支机构":     "6840290 2101",
		"外国（地区）企业":            "7000    2101",
		"外国（地区）公司分支机构":        "7100    2101",
		"外国(地区)无限责任公司分支机构":    "711033052101",
		"外国(地区)有限责任公司分支机构":    "712033052101",
		"外国(地区)股份有限责任公司分支机构":  "713034052101",
		"外国(地区)其他形式公司分支机构":    "719033052101",
		"外国(地区)企业常驻代表机构":      "720033052101",
		"外国(地区)企业在中国境内从事经营活动": "7300    2101",
		//"分公司":         "7310",
		"集团":          "8000    4   ",
		"内资集团":        "8100    4   ",
		"外资集团":        "8500    4   ",
		"其他类型":        "9000   95   ",
		"农民专业合作社":     "91001909155 ",
		"农民专业合作社分支机构": "92001909255 ",
		"个体工商户":       "9500    3   ",
		"自然人":         "9600    5   ",
		"其他":          "9900190 5   ",
	}
	//证件类型
	certTypeDict = map[string]string{
		"中华人民共和国居民身份证": "10",
		"中华人民共和国军官证":   "20",
		"中华人民共和国警官证":   "30",
		"外国(地区)护照":     "40",
		"其他有效身份证件":     "90",
	}
	//币种
	curDict = map[string]string{
		"澳大利亚元": "036",
		"奥地利先令": "040",
		"新加坡元":  "702",
		"人民币":   "156",
		"港币":    "344",
		"港元":    "344",
		"意大利里拉": "380",
		"日元":    "392",
		"韩元":    "410",
		"荷兰盾":   "528",
		"新西兰元":  "554",
		"挪威克朗":  "578",
		"欧元":    "954",
		"瑞典克朗":  "752",
		"瑞士法郎":  "756",
		"英镑":    "826",
		"美元":    "840",
		"比利时法郎": "056",
		"加拿大元":  "124",
		"丹麦克朗":  "208",
		"法国法郎":  "250",
		"德国马克":  "280",
	}
	//证照类型
	blicTypeDict = map[string]string{
		"内资企业法人":            "A1",
		"企业法人营业执照(公司)":      "A1",
		"外资企业法人":            "A2",
		"企业法人营业执照(外资)":      "A2",
		"营业执照(外资)":          "A2",
		"外商投资企业办事机构注册证":     "A2",
		"企业法人营业执照(非公司)":     "B",
		"合伙企业营业执照":          "C",
		"农民专业合作社":           "D",
		"农民专业合作社法人营业执照":     "D",
		"个人独资企业营业执照":        "E",
		"个体工商户营业执照":         "F",
		"非法人企业、个体工商户":       "G",
		"营业执照(分公司、营业单位)":    "G1",
		"合伙企业分支机构营业执照":      "G4",
		"农民专业合作社分支机构营业执照":   "H",
		"企业集团登记证":           "I",
		"外国(地区)企业常驻代表机构登记证": "J",
		"机关、事业、社团法人":        "K",
		"事业法人登记证":           "K1",
		"社团法人登记证":           "K2",
		"机关法人登记证":           "K3",
		"不需要登记的社团法人":        "K4",
		"其他": "Z",
		"1":  "A1",
		"11": "A1",
		"2":  "A2",
		"21": "A2",
		"22": "A2",
		"24": "A2",
		"12": "B",
		"32": "C",
		"8":  "D",
		"81": "D",
		"31": "E",
		"33": "F",
		"3":  "G",
		"14": "G1",
		"36": "G4",
		"82": "H",
		"13": "I",
		"23": "J",
		"6":  "K",
		"61": "K1",
		"62": "K2",
		"63": "K3",
		"66": "K4",
		"9":  "Z",
	}
	//认缴出资方式
	conFormDict = map[string]string{
		"货币":         "1",
		"实物":         "2",
		"知识产权":       "3",
		"知识产权－非专利技术": "3",
		"债权": "4",
		"知识产权－高新技术成果": "3",
		"其他知识产权":      "3",
		"人力资本":        "9",
		"土地使用权":       "6",
		"商标使用权 ":      "9",
		"其他财产 ":       "9",
		"股权":          "7",
		"劳务":          "9",
		"其他":          "9",
	}
	//国籍
	countryDict = map[string]string{
		"阿富汗":      "004",
		"阿尔巴尼亚":    "008",
		"南极洲":      "010",
		"阿尔及利亚":    "012",
		"美属萨摩亚":    "016",
		"安道尔":      "020",
		"安哥拉":      "024",
		"安提瓜和巴布达":  "028",
		"阿塞拜疆":     "031",
		"阿根廷":      "032",
		"澳大利亚":     "036",
		"奥地利":      "040",
		"巴哈马":      "044",
		"巴林":       "048",
		"孟加拉国":     "050",
		"亚美尼亚":     "051",
		"巴巴多斯":     "052",
		"比利时":      "056",
		"百慕大":      "060",
		"不丹":       "064",
		"玻利维亚":     "068",
		"波黑":       "070",
		"博茨瓦那":     "072",
		"布维岛":      "074",
		"巴西":       "076",
		"伯利茨":      "084",
		"英属印度洋领土":  "086",
		"所罗门群岛":    "090",
		"英属维尔京群岛":  "092",
		"文莱":       "096",
		"保加利亚":     "100",
		"缅甸":       "104",
		"布隆迪":      "108",
		"白俄罗斯":     "112",
		"柬埔寨":      "116",
		"喀麦隆":      "120",
		"加拿大":      "124",
		"佛得角":      "132",
		"开曼群岛":     "136",
		"中非":       "140",
		"斯里兰卡":     "144",
		"乍得":       "148",
		"智利":       "152",
		"中国":       "156",
		"中国台湾省":    "158",
		"圣诞岛":      "162",
		"可可(基林)群岛": "166",
		"哥伦比亚":     "170",
		"科摩罗":      "174",
		"马约特":      "175",
		"刚果(布)":    "178",
		"刚果(金)":    "180",
		"库克群岛":     "184",
		"哥斯达黎加":    "188",
		"克罗地亚":     "191",
		"古巴":       "192",
		"塞浦路斯":     "196",
		"捷克":       "203",
		"贝宁":       "204",
		"丹麦":       "208",
		"多米尼克":     "212",
		"多米尼加共和国":  "214",
		"厄瓜多尔":     "218",
		"萨尔瓦多":     "222",
		"赤道几内亚":    "226",
		"埃塞俄比亚":    "231",
		"厄立特里亚":    "232",
		"爱沙尼亚":     "233",
		"法罗":       "234",
		"福克兰群岛(马尔维纳斯)": "238",
		"南乔治亚岛和南桑德韦奇岛": "239",
		"斐济":         "242",
		"芬兰":         "246",
		"法国":         "250",
		"法属圭亚那":      "254",
		"法属波利尼西亚":    "258",
		"法属南部领土":     "260",
		"吉布提":        "262",
		"加蓬":         "266",
		"格鲁吉亚":       "268",
		"冈比亚":        "270",
		"巴勒斯坦":       "275",
		"德国":         "276",
		"加纳":         "288",
		"直布罗陀":       "292",
		"基里巴斯":       "296",
		"希腊":         "300",
		"格陵兰":        "304",
		"格林纳达":       "308",
		"瓜德罗普":       "312",
		"关岛":         "316",
		"危地马拉":       "320",
		"几内亚":        "324",
		"圭亚那":        "328",
		"海地":         "332",
		"赫德岛":        "334",
		"梵帝冈":        "336",
		"洪都拉斯":       "340",
		"中国香港特别行政区":  "344",
		"匈牙利":        "348",
		"冰岛":         "352",
		"印度":         "356",
		"黑山共和国":      "359",
		"印度尼西亚":      "360",
		"伊朗":         "364",
		"伊拉克":        "368",
		"爱尔兰":        "372",
		"以色列":        "376",
		"意大利":        "380",
		"科特迪瓦":       "384",
		"牙买加":        "388",
		"日本":         "392",
		"哈萨克斯坦":      "398",
		"约旦":         "400",
		"肯尼亚":        "404",
		"朝鲜":         "408",
		"韩国":         "410",
		"科威特":        "414",
		"吉尔吉斯坦":      "417",
		"老挝":         "418",
		"黎巴嫩":        "422",
		"莱索托":        "426",
		"拉脱维亚":       "428",
		"利比里亚":       "430",
		"利比亚":        "434",
		"列支敦士登":      "438",
		"立陶宛":        "440",
		"卢森堡":        "442",
		"中国澳门特别行政区":  "446",
		"马达加斯加":      "450",
		"马拉维":        "454",
		"马来西亚":       "458",
		"马尔代夫":       "462",
		"马里":         "466",
		"马耳他":        "470",
		"马提尼克":       "474",
		"毛里塔尼亚":      "478",
		"毛里求斯":       "480",
		"墨西哥":        "484",
		"摩纳哥":        "492",
		"蒙古":         "496",
		"摩尔多瓦":       "498",
		"蒙特塞拉特":      "500",
		"摩洛哥":        "504",
		"莫桑比克":       "508",
		"阿曼":         "512",
		"纳米比亚":       "516",
		"瑙鲁":         "520",
		"尼泊尔":        "524",
		"荷兰":         "528",
		"荷属安的列斯":     "530",
		"阿鲁巴":        "533",
		"新喀里多尼亚":     "540",
		"瓦努阿图":       "548",
		"新西兰":        "554",
		"尼加拉瓜":       "558",
		"尼日尔":        "562",
		"尼日利亚":       "566",
		"纽埃":         "570",
		"诺福克岛":       "574",
		"挪威":         "578",
		"北马里亚纳":      "580",
		"密克罗尼西亚联邦":   "583",
		"马绍尔群岛共和国":   "584",
		"帕劳":         "585",
		"巴基斯坦":       "586",
		"巴拿马":        "591",
		"巴布亚新几内亚":    "598",
		"巴拉圭":        "600",
		"秘鲁":         "604",
		"菲律宾":        "608",
		"皮特凯恩群岛":     "612",
		"波兰":         "616",
		"葡萄牙":        "620",
		"几内亚比绍":      "624",
		"东帝汶":        "626",
		"波多黎各":       "630",
		"卡塔尔":        "634",
		"留尼汪":        "638",
		"罗马尼亚":       "642",
		"俄罗斯":        "643",
		"卢旺达":        "646",
		"圣赫勒拉":       "654",
		"圣基茨和尼维斯":    "659",
		"安圭拉":        "660",
		"圣卢西亚":       "662",
		"圣皮埃尔和密克隆":   "666",
		"圣文森特和格林纳丁斯": "670",
		"圣马力诺":       "674",
		"圣多美和普林西比":   "678",
		"沙特阿拉伯":      "682",
		"塞内加尔":       "686",
		"塞舌尔":        "690",
		"塞拉利昂":       "694",
		"新加坡":        "702",
		"斯洛伐克":       "703",
		"越南":         "704",
		"斯洛文尼亚":      "705",
		"索马里":        "706",
		"南非":         "710",
		"津巴布维":       "716",
		"西班牙":        "724",
		"西撒哈垃":       "732",
		"苏丹":         "736",
		"苏里南":        "740",
		"斯瓦尔巴岛和扬马延岛": "744",
		"斯威士兰":       "748",
		"瑞典":         "752",
		"瑞士":         "756",
		"叙利亚":        "760",
		"塔吉克斯坦":      "762",
		"泰国":         "764",
		"多哥":         "768",
		"托克劳":        "772",
		"汤加":         "776",
		"特立尼达和多巴哥":   "780",
		"阿拉伯联合酋长国":   "784",
		"突尼斯":        "788",
		"土耳其":        "792",
		"土库曼斯坦":      "795",
		"特克斯和凯科斯群岛":  "796",
		"图瓦卢":        "798",
		"乌干达":        "800",
		"乌克兰":        "804",
		"马其顿共和国":     "807",
		"埃及":         "818",
		"英国":         "826",
		"坦桑尼亚":       "834",
		"美国":         "840",
		"美属维尔京群岛":    "850",
		"布基那法索":      "854",
		"乌拉圭":        "858",
		"乌兹别克斯坦":     "860",
		"委内瑞拉":       "862",
		"瓦利斯和富图纳群岛":  "876",
		"萨摩亚":        "882",
		"阿拉伯也门共和国":   "887",
		"也门":         "887",
		"南斯拉夫":       "891",
		"赞比亚":        "894",
		"其他国家地区":     "   ",
	}
	//变更事项
	altitemDict = map[string]string{
		"名称变更":                  "01",
		"住所变更":                  "02",
		"法定代表人变更":               "03",
		"企业类型变更":                "04",
		"注册资本(金)变更":             "20",
		"许可经营项目变更":              "10",
		"一般经营项目变更":              "10",
		"经营期限(营业期限)变更":          "05",
		"经营方式变更":                "05",
		"投资人(股权)变更":             "30",
		"行业代码变更":                "10",
		"经营范围变更":                "10",
		"经营场所变更":                "10",
		"登记机关变更":                "06",
		"出资方式变更":                "22",
		"出资比例变更":                "22",
		"出资日期变更":                "22",
		"出资额变更":                 "22",
		"投资总额变更":                "22",
		"负责人变更":                 "03",
		"首席代表变更":                "03",
		"驻在地址变更":                "02",
		"驻在期限变更":                "02",
		"业务范围变更":                "10",
		"资金数额变更":                "22",
		"实收资本变更":                "22",
		"执行合伙企业事务的合伙人":          "03",
		"投资者名称(姓名)变更":           "03",
		"增加分支机构变更":              "72",
		"撤销分支机构变更":              "72",
		"企业产权合同变更":              "03",
		"投资人居所变更":               "02",
		"其他变更":                  "99",
		"董事备案":                  "70",
		"监事备案":                  "70",
		"经理备案":                  "70",
		"章程备案":                  "71",
		"章程修正案备案":               "71",
		"股权质押备案":                "71",
		"实缴出资备案":                "71",
		"出资期限变更备案":              "71",
		"清算组成员备案":               "73",
		"企业印章、印模备案":             "99",
		"分公司增加备案":               "72",
		"分公司变更备案":               "72",
		"分公司注销备案":               "72",
		"境外股东发起人的境内法律文件送达接受人备案": "75",
		"工商登记联络员备案":             "78",
		"其他机构备案":                "99",
		"其他事项备案":                "99",
	}
	//登记机关,处理地和行政区划
	areacodeDict = map[string]string{
		"上海市黄浦区":  "310101",
		"上海市徐汇区":  "310104",
		"上海市长宁区":  "310105",
		"上海市静安区":  "310106",
		"上海市普陀区":  "310107",
		"上海市闸北区":  "310106",
		"上海市虹口区":  "310109",
		"上海市杨浦区":  "310110",
		"上海市闵行区":  "310112",
		"上海市宝山区":  "310113",
		"上海市嘉定区":  "310114",
		"上海市浦东新区": "310115",
		"上海市金山区":  "310116",
		"上海市松江区":  "310117",
		"上海市青浦区":  "310118",
		"上海市奉贤区":  "310120",
		"上海市崇明区":  "310151",
		"上海市崇明县":  "310151",
		"黄浦区":     "310101",
		"徐汇区":     "310104",
		"长宁区":     "310105",
		"静安区":     "310106",
		"普陀区":     "310107",
		"闸北区":     "310106",
		"虹口区":     "310109",
		"杨浦区":     "310110",
		"闵行区":     "310112",
		"宝山区":     "310113",
		"嘉定区":     "310114",
		"浦东新区":    "310115",
		"金山区":     "310116",
		"松江区":     "310117",
		"青浦区":     "310118",
		"奉贤区":     "310120",
		"崇明区":     "310151",
		"崇明县":     "310151",
		"黄浦":      "310101",
		"徐汇":      "310104",
		"长宁":      "310105",
		"静安":      "310106",
		"普陀":      "310107",
		"闸北":      "310106",
		"虹口":      "310109",
		"杨浦":      "310110",
		"闵行":      "310112",
		"宝山":      "310113",
		"嘉定":      "310114",
		"浦东":      "310115",
		"金山":      "310116",
		"松江":      "310117",
		"青浦":      "310118",
		"奉贤":      "310120",
		"崇明":      "310151",
		"市属":      "310000",
	}
	//投资类型
	tzTypeDict = map[string]string{
		"内资企业法人":            "11",
		"企业法人营业执照(公司)":      "11",
		"外资企业法人":            "15",
		"企业法人营业执照(外资)":      "15",
		"营业执照(外资)":          "15",
		"外商投资企业办事机构注册证":     "90",
		"企业法人营业执照(非公司)":     "11",
		"合伙企业营业执照":          "40",
		"农民专业合作社":           "90",
		"农民专业合作社法人营业执照":     "90",
		"个人独资企业营业执照":        "50",
		"个体工商户营业执照":         "90",
		"非法人企业、个体工商户":       "90",
		"营业执照(分公司、营业单位)":    "90",
		"合伙企业分支机构营业执照":      "90",
		"农民专业合作社分支机构营业执照":   "90",
		"企业集团登记证":           "90",
		"外国(地区)企业常驻代表机构登记证": "90",
		"机关、事业、社团法人":        "90",
		"事业法人登记证":           "12",
		"社团法人登记证":           "13",
		"机关法人登记证":           "14",
		"不需要登记的社团法人":        "13",
		"其他": "90",
	}
)

func convert(dict map[string]string, val string) string {
	vald:=strings.Replace(strings.Replace(val,"（","(",-1),"）",")",-1) 
	val, ok := dict[strings.TrimSpace(vald)]
	if ok {
		return val
	}
	return ""
}

//根据公司类型，依次得到登记注册类型,控股情况，单位类型，机构类型，执行会计标准
func transform(dict map[string]string, val string) (string, string, string, string, string, string) {
	vald:=strings.Replace(strings.Replace(val,"（","(",-1),"）",")",-1) 
	vals, ok := dict[strings.TrimSpace(vald)]
	if ok {
		rs := []rune(vals)
		return string(rs[0:4]), string(rs[4:7]), string(rs[7:8]), string(rs[8:9]), string(rs[9:11]), string(rs[11:])
	}
	return "", "", "", "", "", ""
}
