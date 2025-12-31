package model

type ExerciseList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ExerciseResponse struct {
	Total     int            `json:"total"`
	Limit     int            `json:"limit"`
	Offset    int            `json:"offset"`
	Count     int            `json:"count"`
	Exercises []ExerciseList `json:"results"`
}

type ExerciseDetail struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	PrimaryMuscles []string        `json:"primaryMuscles"`
	Category       string          `json:"category"`
	Force          string          `json:"force"`
	Grips          []string        `json:"grips"`
	Mechanic       string          `json:"mechanic"`
	Difficulty     string          `json:"difficulty"`
	Steps          []string        `json:"steps"`
	Videos         []ExerciseVideo `json:"videos"`
	VideoCount     int             `json:"video_count"`
	StepCount      int             `json:"step_count"`
}

type ExerciseVideo struct {
	Url     string `json:"url"`
	Angle   string `json:"angle"`
	Gender  string `json:"gender"`
	OgImage string `json:"og_image"`
}
