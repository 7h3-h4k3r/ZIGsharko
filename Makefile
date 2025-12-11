
RESET   = \033[0m
RED     = \033[0;31m
GREEN   = \033[0;32m
YELLOW  = \033[0;33m
BLUE    = \033[0;34m
CYAN    = \033[0;36m
WHITE   = \033[0;37m
BOLD    = \033[1m

package_install's:
	@echo "$(CYAN)Installing Package . . . $(RESET)"
	go get github.com/go-chi/chi/v5
	go get github.com/go-chi/chi/v5/middleware
	go get github.com/go-playground/validator/v10 