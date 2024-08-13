package pushbell

type Urgency string

const (
	UrgencyVeryLow Urgency = "very-low" // Device State - On power and Wi-Fi
	UrgencyLow     Urgency = "low"      // Device State - On either power or Wi-Fi
	UrgencyNormal  Urgency = "normal"   // Device State - On neither power nor Wi-Fi
	UrgencyHigh    Urgency = "high"     // Device State - Low battery
)
