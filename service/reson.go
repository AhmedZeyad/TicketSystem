package service

import (
	"TicketSystem/engine"
	"errors"

	// "fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Reason struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Status bool   `json:"status" db:"status"`
	// done  add created and updated and deleted at ,by as feald
	CreatedAt time.Time `db:"created_at" json:"-"`
	CreatedBy int       `db:"created_by" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
	UpdatedBy int       `db:"updated_by" json:"-"`
	DeletedAt time.Time `db:"deleted_at" json:"-"`
	DeletedBy int       `db:"deleted_by" json:"-"`
}

type SubReason struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	ReasonId int    `json:"reasonId" db:"reasonId"`
}

// todo updae select  conditon
func SelectReasons() ([]Reason, error) {
	var reasons []Reason
	err := engine.DB.Select(&reasons, "SELECT id, name,  status FROM reasons where deleted_at is not null")
	if err != nil {
		return nil, err
	}
	if len(reasons) == 0 {
		return reasons, errors.New("no reasons found")
	}
	return reasons, nil
}
func InsertReason(reason Reason) (Reason, error) {

	resoult, err := engine.DB.Exec("INSERT INTO reasons( name ,updated_by,created_by)   VALUES (?,?,?)   ", reason.Name, 1, 1)
	if err != nil {

		return Reason{}, err
	}
	rowCount, err := resoult.RowsAffected()
	if err != nil {
		return Reason{}, err
	}
	if rowCount == 0 {
		return Reason{}, errors.New("failed to insert reason")
	}
	reason.ID = int(rowCount)
	return reason, nil
}

func SelectReasonById(id int) (Reason, error) {
	var reason Reason
	err := engine.DB.Get(&reason, "SELECT id, name, status  FROM reasons WHERE id = ? and WHERE created_ad is not null", id)
	if err != nil {
		return reason, err
	}

	return reason, nil
}

// todo add update and dlete

// todo add reason updateDB

func UpdateReason(reason Reason) (state bool, err error) {
	resoult, err := engine.DB.Exec("UPDATE reasons SET name = ?, status = ?  WHERE id = ?", reason.Name, reason.Status, reason.ID)
	if err != nil {
		return false, err
	}
	affecRow, err := resoult.RowsAffected()
	if err != nil {
		return false, err
	}
	if affecRow == 0 {
		return true, errors.New("failed to update reason")
	}
	return true, nil
}

// done add reason deleteDB
func DeleteReason(reasonID, userId int) (bool, error) {
	resoult, err := engine.DB.Exec("UPDATE reasons SET deleted_at = ?  deleted_by = ? WHERE id = ?", time.Now(), userId, reasonID)
	if err != nil {
		return false, err
	}
	affectedRow, err := resoult.RowsAffected()
	if err != nil {
		return false, err

	}
	if affectedRow == 0 {
		return false, errors.New("failed to delete reason")
	}
	return true, nil
}

// done update query
func SelecteSubreasonByReasonId(id int) ([]SubReason, error) {
	var subreasons []SubReason
	err := engine.DB.Select(&subreasons, `SELECT id ,name ,reasonId  FROM subReasons WHERE reasonId =?`, id)
	if err != nil {
		return nil, err
	}
	if len(subreasons) == 0 {
		return subreasons, errors.New("no sub reasons found")
	}
	return subreasons, nil
}
func (sub SubReason) InsertSubReason() error {
	reoult, err := engine.DB.Exec(`INSERT INTO subReasons (name,reasonId)  VALUES (?,?) `, sub.Name, sub.ReasonId)
	if err != nil {
		return err
	}
	affecRow, err := reoult.RowsAffected()
	if err != nil {
		return err
	}
	if affecRow == 0 {
		return errors.New(" 0 row  insert to sub reason")
	}

	return nil

}

// todo update and dlete
// done update

func (subReason SubReason) SelectSubReasonById() (SubReason, error) {
	resoult, err := engine.DB.Exec("update subReasons SET name = ?  WHERE id = ?", subReason.Name, subReason.ID)
	if err != nil {
		return SubReason{}, err
	}

	affecRow, err := resoult.RowsAffected()
	if err != nil {
		return SubReason{}, err
	}
	if affecRow == 0 {
		return SubReason{}, errors.New("failed to update sub reason")
	}
	return subReason, nil

}

// todo delete
func DeleteSubReasonById(id int) (bool, error) {
	resoult, err := engine.DB.Exec("DELETE subReasons   WHERE id = ?", time.Now(), id)
	if err != nil {
		return false, err
	}

	affecRow, err := resoult.RowsAffected()
	if err != nil {
		return false, err
	}
	if affecRow == 0 {
		return false, errors.New("failed to delete sub reason")
	}
	return true, nil
}

// api

func TicketReasonRoutes(api *gin.RouterGroup) {
	// reason group
	reasonRoutes := api.Group("reason")

	// get all
	reasonRoutes.GET("/", func(c *gin.Context) {
		r, err := SelectReasons()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  err.Error(),
				"message": "can't get reasons",
			})
			return
		}

		// if the there is no reason
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"reasons": r,
			},
		})
	})
	// get by id
	reasonRoutes.GET("/:id", func(c *gin.Context) {
		strID := c.Param("id")
		id, err := strconv.Atoi(strID)
		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  err.Error(),
				"message": "Not valid id",
			})
			return
		}
		reason, err := SelectReasonById(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  err.Error(),
				"message": "can't get reason",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"reason": reason,
			},
		})

	})
	// add
	reasonRoutes.POST("/", func(c *gin.Context) {
		var reason Reason
		if err := c.ShouldBindJSON(&reason); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  err.Error(),
				"message": "not vallid data",
			})
			return
		}
		_, err := InsertReason(reason)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  err.Error(),
				"message": "can't add reason",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"reason": reason,
			},
		})
	})
	// sub reason group
	subreasonRoutes := api.Group("/subreason")
	// get by reason id
	subreasonRoutes.GET("/:id", func(c *gin.Context) {
		strId := c.Param("id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  err.Error(),
				"message": "Not valid id",
			})
			return
		}
		sub_reason, err := SelecteSubreasonByReasonId(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  err.Error(),
				"message": "can't get sub_reason",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"reason": sub_reason,
			},
		})
	})
	// add
	subreasonRoutes.POST("/", func(c *gin.Context) {
		var sub_reason SubReason
		if err := c.BindJSON(&sub_reason); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  err.Error(),
				"message": "not vallid data",
			})
			return
		}
		err := sub_reason.InsertSubReason()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  err.Error(),
				"message": "can't add sub_reason",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"sub_reason": sub_reason,
			},
		})
	})
	// todo update and dlete

	// todo update

	// todo delete

}
