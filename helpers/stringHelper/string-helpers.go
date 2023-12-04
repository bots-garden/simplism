package stringHelper

 func GetTheBooleanValueOf(value string) bool {
    switch {
    case value == "true":
        return true
    case value == "false":
        return false
    default:
        return false
    }
}

func GetTheStringValueOf(value bool) string {
    if value {
        return "true"
    } else {
        return "false"
    }
}