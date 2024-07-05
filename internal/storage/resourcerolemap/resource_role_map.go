package resourcerolemap

import (
	"fmt"

	"github.com/ardihikaru/graphql-example-part-1/pkg/mysqldb"
)

const (
	table = "resource_role_map"
)

type Storage mysqldb.Storage
type statement mysqldb.DbQuery

// ResourceRoleMap is the table property
type ResourceRoleMap struct {
	Id       int    `db:"id"`
	Resource string `db:"resource"`
	Role     string `db:"role"`
}

// GetObjListOwner fetches list of objects who has access to this resource
//
//	 it purposely uses a list to prepare any future update from the table,
//		e.g. updates on the table schema which leads to a multiple records as the result
func (s *Storage) GetObjListOwner(resource, role string) ([]string, error) {
	var err error
	var st statement

	st.Q = fmt.Sprintf(`
		%s
		SELECT *
		FROM %s
		WHERE resource = '%s' AND role = '%s'
	`, st.Q, table, resource, role)

	var rsltMulti []ResourceRoleMap
	err = s.Db.Select(&rsltMulti, st.Q)
	if err != nil {
		s.Log.Error(fmt.Sprintf("got error in query: %s", err.Error()))
		return nil, err
	}

	var objList []string
	resourceMap := map[string]int{} // simply to flag whether key resource exists or not in this map

	for k := range rsltMulti {
		// process only if the key does not exist yet
		if _, exists := resourceMap[rsltMulti[k].Resource]; !exists {
			resourceMap[rsltMulti[k].Resource] = rsltMulti[k].Id

			// appends into the final result
			objList = append(objList, rsltMulti[k].Resource)
		}

	}

	return objList, nil
}
