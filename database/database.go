package database

import (
	"database/sql"
	"fmt"

	"github.com/go-ole/go-ole"
)

type Database struct {
	Path string

	db *sql.DB
}

func (d *Database) Connect() error {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	provider := "Microsoft.ACE.OLEDB.12.0"
	var err error
	d.db, err = sql.Open("adodb",
		fmt.Sprintf("Provider=%s;Data Source=%s;", provider, d.Path),
	)

	if err != nil {
		return fmt.Errorf("unable to open database: %s", err)
	}

	err = d.db.Ping()
	if err != nil {
		return fmt.Errorf("unable to ping database: %T: error \"%s\"", err, err)
	}

	return nil
}

func (d *Database) ActiveOrders() ([]Order, error) {
	rows, err := d.db.Query(`
	select 
		tblOrder.*,tblStates.* 
	from 
		tblOrder,tblStates 
	where 
		tblOrder.State=tblStates.State`)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %s", err)
	}

	orders := make([]Order, 0)

	for rows.Next() {
		order := Order{State: State{}, Positions: make([]OrderPos, 0)}
		err := rows.Scan(
			&order.ONo,
			&order.PlanedStart,
			&order.PlanedEnd,
			&order.Start,
			&order.End,
			&order.CNo,
			&order.StateID,
			&order.Enabled,
			&order.Release,
			&order.State.State,
			&order.State.Description,
			&order.State.Short,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %s", err)
		}

		// further populate Positions
		pos, err := d.db.Query(`
			select 
				tblOrderPos.*, tblResource.*
			from 
				tblOrderPos, tblResource
			where
				tblOrderPos.ONo = ?
			and
				tblOrderPos.ResourceID = tblResource.ResourceID
		`, order.ONo)

		if err != nil {
			return nil, fmt.Errorf("unable to fetch positions, resources and carrier: %s", err)
		}

		for pos.Next() {
			p := OrderPos{Resource: Resource{}}
			err = pos.Scan(
				&p.ONo,
				&p.OPos,
				&p.PlanedStart,
				&p.PlanedEnd,
				&p.Start,
				&p.End,
				&p.WPNo,
				&p.StepNo,
				&p.MainOPos,
				&p.State,
				&p.ResourceID,
				&p.OpNo,
				&p.WONo,
				&p.PNo,
				&p.subOrderBlocked,
				&p.Error,

				&p.Resource.ResourceID,
				&p.Resource.ResourceName,
				&p.Resource.Description,
				&p.Resource.PlcType,
				&p.Resource.IP,
				&p.Resource.Picture,
				&p.Resource.ParallelProcessing,
				&p.Resource.Automatic,
				&p.Resource.WebPage,
				&p.Resource.DefaultBrowser,
				&p.Resource.TopologyType,
			)

			if err != nil {
				return nil, fmt.Errorf("unable to scan orderpos: %w", err)
			}

			p.Carrier, err = d.carrier(order.ONo, p.OPos)
			if err != nil {
				return nil, fmt.Errorf("unable to fetch carrier: %w", err)
			}

			order.Positions = append(order.Positions, p)

		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (d *Database) carrier(no, pos int) (*Carrier, error) {
	c := &Carrier{}
	rows, err := d.db.Query(`
		select 
			* 
		from 
			tblCarrier
		where
			ONo = ?
		and
			OPos = ?
			`,
		no, pos)
	if err != nil {
		return nil, fmt.Errorf("could not query carrier: %s", err)
	}

	if !rows.Next() {
		return nil, nil
	}

	err = rows.Scan(
		&c.CarrierID,
		&c.CarrierTypeID,
		&c.ONo,
		&c.OPos,
		&c.PNo,
		&c.PNoGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("could not scan carrier: %s", err)
	}
	return c, nil
}

func (d *Database) PreviousOrders(ONos ...string) ([]FinOrder, error) {
	orders := make([]FinOrder, 0)
	for _, ono := range ONos {
		rows, err := d.db.Query("select * from tblFinOrder where ONo = ?", ono)
		if err != nil {
			return nil, fmt.Errorf("unable to execute query: %w", err)
		}

		for rows.Next() {
			order := FinOrder{}
			err := rows.Scan(
				&order.ONo,
				&order.PlanedStart,
				&order.PlanedEnd,
				&order.Start,
				&order.End,
				&order.CNo,
				&order.State,
				&order.Enabled,
				&order.Release,
			)
			if err != nil {
				return nil, fmt.Errorf("unable to scan row: %s", err)
			}
			orders = append(orders, order)
		}
	}

	return orders, nil
}
