package handlers

//func CreateBackup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	dbParam := r.URL.Query().Get("db")
//	if dbParam != "" {
//		list := strings.Split(dbParam, ",")
//		if err := manager.Backup(list); err != nil {
//			_, _ = fmt.Fprint(w, err.Error())
//			return
//		}
//	} else {
//		if err := manager.BackupAll(); err != nil {
//			_, _ = fmt.Fprint(w, err.Error())
//			return
//		}
//	}
//
//	_, _ = fmt.Fprint(w, "OK")
//}
//
//func RestoreBackup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	backupFile := ps.ByName("name")
//	if backupFile == "latest" {
//		RestoreLastBackup(w, r, ps)
//		return
//	}
//
//	dbParam := r.URL.Query().Get("db")
//	if dbParam != "" {
//		list := strings.Split(dbParam, ",")
//		if err := manager.Restore(list, backupFile); err != nil {
//			_, _ = fmt.Fprint(w, err.Error())
//			return
//		}
//	} else {
//		if err := manager.RestoreAll(backupFile); err != nil {
//			_, _ = fmt.Fprint(w, err.Error())
//			return
//		}
//	}
//
//	_, _ = fmt.Fprint(w, "OK")
//}
//
//func RestoreLastBackup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	bucket := manager.Storage
//	list, err := bucket.List()
//	if err != nil {
//		_, _ = fmt.Fprint(w, err.Error())
//		return
//	}
//
//	if list == nil {
//		_, _ = fmt.Fprint(w, "List is nil")
//		return
//	}
//
//	backupFile := list[len(list) - 1]
//	if err := manager.RestoreAll(backupFile); err != nil {
//		_, _ = fmt.Fprint(w, err.Error())
//		return
//	}
//
//	_, _ = fmt.Fprint(w, "OK")
//}
//
//func ListBackup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	bucket := manager.Storage
//	list, err := bucket.List()
//	if err != nil {
//		_, _ = fmt.Fprint(w, err.Error())
//		return
//	}
//
//	js, err := json.Marshal(list)
//	if err != nil {
//		_, _ = fmt.Fprint(w, err.Error())
//		return
//	}
//
//	_, _ = fmt.Fprint(w, string(js))
//}
//
//func DeleteBackup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	bucket := manager.Storage
//	backupFile := ps.ByName("name")
//	if err := bucket.Delete(backupFile); err != nil {
//		_, _ = fmt.Fprint(w, err.Error())
//		return
//	}
//
//	_, _ = fmt.Fprint(w, "OK")
//}
