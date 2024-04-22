# Migration

If migration is to be performed with dynmgrm, you must use the `dynmgrm` tag.

```.go
package main

type User struct {
	ProjectID string `gorm:"primaryKey, dynmgrm:"pk"`
	ID        string `gorm:"primaryKey, dynmgrm:"sk;gsi:id_name-index,pk"`
	Name      string `dynmgrm:"gsi:id_name-index,sk;lsi:name-index"`
	Note      string `dynmgrm:"non-projective:[id_name-index,name-index]"`
}
```

## Fields Tags

| Tag Name       | Format                      | Description                                                                                                              |
|----------------|-----------------------------|--------------------------------------------------------------------------------------------------------------------------|
| pk             | -                           | field will be the PK attribute of the table.                                                                             |
| sk             | -                           | field will be the SK attribute of the table.                                                                             |
| gsi            | :\<index name\>,<pk \| sk\> | field will be the PK or SK attribute of the specified Global Secondary Index.<br/>It works with `Migrator.CreateIndex()` |
| lsi            | :\<index name\>             | field will be the SK attribute of the specified Local Secondary Index.                                                   |
| non-projective | :[(,)\<index name\>]        | exclude from projection at enumerated Index.                                                                             |
