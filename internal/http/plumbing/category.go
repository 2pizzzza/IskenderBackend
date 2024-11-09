package plumbing

//func (s *Server) CreateCategory(w http.ResponseWriter, r *http.Request) {
//	var req schemas.CreateCategoryRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
//		return
//	}
//	defer r.Body.Close()
//
//	s.log.Info("Creating category", slog.String("name", req.Name))
//
//	res, err := s.service.CreateCategory(r.Context(), &req)
//	if err != nil {
//		s.log.Error("Failed to create category", sl.Err(err))
//		http.Error(w, "Failed to create category", http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(res); err != nil {
//		s.log.Error("Failed to encode response", sl.Err(err))
//		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
//		return
//	}
//
//	s.log.Info("Category created successfully", slog.String("name", res.Name))
//}
//
//func (s *Server) GetAllCategories(w http.ResponseWriter, r *http.Request) {
//
//	res, err := s.service.GetAllCategory(r.Context())
//	if err != nil {
//		s.log.Error("Failed get all categories", sl.Err(err))
//		http.Error(w, "Failed get all categories", http.StatusInternalServerError)
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(res); err != nil {
//		s.log.Error("Failed to encode response", sl.Err(err))
//		http.Error(w, "Failed encode response", http.StatusInternalServerError)
//		return
//	}
//}
//
//func (s *Server) GetCategoryById(w http.ResponseWriter, r *http.Request) {
//	idStr := r.URL.Query().Get("id")
//	if idStr == "" {
//		http.Error(w, "Missing category Id", http.StatusBadRequest)
//		return
//	}
//
//	categoryId, err := strconv.Atoi(idStr)
//	if err != nil {
//		http.Error(w, "Invalid item ID", http.StatusBadRequest)
//		return
//	}
//
//	s.log.Info("Fetching category by Item", slog.Int("category_id", categoryId))
//
//	req := &schemas.CategoryByIdRequest{Id: categoryId}
//
//	category, err := s.service.GetCategoryById(r.Context(), req)
//	if err != nil {
//		if errors.Is(err, schemas.ErrItemNotFound) {
//			s.log.Error("Category not found", sl.Err(err))
//			http.Error(w, "Category not found", http.StatusNotFound)
//			return
//		}
//		s.log.Error("Failed to get category by ID", sl.Err(err))
//		http.Error(w, "Failed to get category by ID", http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(category); err != nil {
//		s.log.Error("Failed encode response", sl.Err(err))
//		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
//		return
//	}
//}
//
//func (s *Server) UpdateCategory(w http.ResponseWriter, r *http.Request) {
//	var req schemas.UpdateCategoryRequest
//
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
//		return
//	}
//
//	defer r.Body.Close()
//	s.log.Info("Update category", slog.String("name", req.NewName))
//
//	err := s.service.UpdateCategory(r.Context(), &req)
//	if err != nil {
//		if errors.Is(err, schemas.ErrItemNotFound) {
//			s.log.Error("Category not found", sl.Err(err))
//			http.Error(w, "Category not found", http.StatusNotFound)
//			return
//		}
//		s.log.Error("Failed to update category by ID", sl.Err(err))
//		http.Error(w, "Failed to update category by ID", http.StatusInternalServerError)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	resp := map[string]string{"message": "Category updated successfully"}
//	if err := json.NewEncoder(w).Encode(resp); err != nil {
//		s.log.Error("Failed to encode response", sl.Err(err))
//		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
//		return
//	}
//
//	s.log.Info("Category updated successfully", slog.Int("category_id", req.Id))
//}
//
//func (s *Server) RemoveCategory(w http.ResponseWriter, r *http.Request) {
//	var req schemas.CategoryByIdRequest
//
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
//		return
//	}
//
//	defer r.Body.Close()
//	s.log.Info("Remove category", slog.Int("name", req.Id))
//
//	err := s.service.RemoveCategory(r.Context(), &req)
//	if err != nil {
//		if errors.Is(err, schemas.ErrItemNotFound) {
//			s.log.Error("Category not found", sl.Err(err))
//			http.Error(w, "Category not found", http.StatusNotFound)
//			return
//		}
//		s.log.Error("Failed to remove category by ID", sl.Err(err))
//		http.Error(w, "Failed to remove category by ID", http.StatusInternalServerError)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	resp := map[string]string{"message": "Category remove successfully"}
//	if err := json.NewEncoder(w).Encode(resp); err != nil {
//		s.log.Error("Failed to encode response", sl.Err(err))
//		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
//		return
//	}
//
//	s.log.Info("Category remove successfully", slog.Int("category_id", req.Id))
//}