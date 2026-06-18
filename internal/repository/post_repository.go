package repository

import (

)

type PostRepository interface {
	Create()
	GetByID()
	List()
	Update()
}