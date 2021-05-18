package utility

var (
	SUCCESS = "SM001"
	SYSTEMERROR = "SM002"
	INPUTERROR = "SM003"

	VALIDATIONSUCCESS = "VAL001"
	VALIDATIONFAILED = "VAL002"
)

func GetCodeMsg(code string) string {
	switch code {
	case "SM001":
		return "Request processed successfully"
	case "SM002":
		return "Request process failed. Please try again."
	case "SM003":
		return "Request failed for some input."
	case "VAL001":
		return "Validation was successful."
	case "VAL002":
		return "Validation failed. please check your input."
	}
	return "Unknown request"
}