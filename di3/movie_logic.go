package main

type MovieLogic struct {
	*MovieModel
}

func NewMovieLogic() *MovieLogic {
	return &MovieLogic{}
}
