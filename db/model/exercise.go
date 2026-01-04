package model

// ExerciseList is the list of exercises returned by the MuscleWiki API
type ExerciseList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ExerciseResponse is the response from the MuscleWiki API for a list of exercises
type ExerciseResponse struct {
	Total     int            `json:"total"`
	Limit     int            `json:"limit"`
	Offset    int            `json:"offset"`
	Count     int            `json:"count"`
	Exercises []ExerciseList `json:"results"`
}

// ExerciseDetail is the detail of an exercise returned by the MuscleWiki API
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

// ExerciseVideo is the video of an exercise returned by the MuscleWiki API
type ExerciseVideo struct {
	Url     string `json:"url"`
	Angle   string `json:"angle"`
	Gender  string `json:"gender"`
	OgImage string `json:"og_image"`
}
