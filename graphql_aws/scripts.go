package graphql_aws

var (
	SelectString = "SELECT  %s FROM %s LIMIT %d OFFSET %d"

	CreateTablePrefix = `
				  mutation Create%s(
					$input: Create%sInput!
				  ) {
					create%s(input: $input) {organizationID  
	`
	CreateTableSuffix = `
					}
				  }
	`

	UpdateTablePrefix = `
				  mutation Update%s(
					$input: Update%sInput!
				  ) {
					update%s(input: $input) {
    `
	UpdageTableSuffix = `
					}
				  }
	`
	QueryTableSuffix = `
			  query Get%s($id: ID!) {
				get%s(id: $id) {
				  id
			  }
			}
    `
	DeleteTableSuffix = `
			 mutation Delete%s(
				$input: Delete%sInput!
				$condition: Model%sConditionInput
			  ) {
				delete%s(input: $input, condition: $condition) {
				  id
				}
			  }
		`
	ListPrefix = `
			  query List%ss(
				$filter: Model%sFilterInput
				$limit: Int
				$nextToken: String
			  ) {
				list%ss(filter: $filter, limit: $limit, nextToken: $nextToken) {
				  items {
					id
				  }
				  nextToken
				}
			  }
    `
	GetOrganization = `
		  query GetOrganization($id: ID!) {
			getOrganization(id: $id) {
			  id
			  name
			  syncSFObject
			  sfLoginUrl
			  sfUsername
			  sfPassword
			  sfSecretToken
			  }
			}
    `
	ListSFObjects = `
			  query ListSFObjects(
				$limit: Int
				$filter: ModelSFObjectFilterInput
			  ) {
				listSFObjects(filter: $filter, limit: $limit) {
				  items {
					id
					name
					organizationID
					name_jp
					sfObjectName
					fields {
					  order
					  name
					  name_jp
					  apiName
					}
				  }
				}
			  }
    `
	ListOrganizations = `
	  query ListOrganizations {
		listOrganizations {
		  items {
			id
			name
			sfLoginUrl
			sfUsername
			sfPassword
			sfSecretToken
			syncSFObject
			sfObjects {
			  items {
				id
				name
				sfObjectName
				fields {
				  name
				  apiName
				}
			  }
			}
		  }
		}
	  }
    `
	gListMetaTables = `
		  query ListMetaTables(
			$filter: ModelMetaTableFilterInput
		  ) {
			listMetaTables(filter: $filter) {
			  items {
				id
				name
				name_jp
				fields {
				  order
				  name
				  name_jp
				}
				createdAt
				updatedAt
			  }
			}
		  }
    `

	ListMetaTables = `
				  query ListMetaTables(
					$filter: ModelMetaTableFilterInput
				  ) {
					listMetaTables(filter: $filter) {
					  items {
						id
						name
						name_jp
						fields {
						  order
						  name
						  name_jp
						}
						createdAt
						updatedAt
					  }
					}
				  }
    `
)
