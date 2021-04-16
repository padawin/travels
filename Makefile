ifeq ($(SOURCE_DIR),)
$(error "SOURCE_DIR missing")
endif

ifeq ($(DEST_DIR),)
$(error "DEST_DIR missing")
endif

ifeq ($(JSON_FILE),)
$(error "JSON_FILE missing")
endif

images := $(shell find $(SOURCE_DIR) -iname '*.jpg' | sed -e "s/ /\\\ /g" | sort)
dest_images := $(subst $(SOURCE_DIR),$(DEST_DIR)/1024x768x0,$(images))\
			   $(subst $(SOURCE_DIR),$(DEST_DIR)/100x100x1,$(images))\
			   $(subst $(SOURCE_DIR),$(DEST_DIR)/118x133x1,$(images))\
			   $(subst $(SOURCE_DIR),$(DEST_DIR)/250x250x1,$(images))
dest_images := $(subst :,\:,$(dest_images))

.PHONY: all
all: $(dest_images)
	./bin/generate.py json '$(JSON_FILE)' $(SOURCE_DIR)

$(DEST_DIR)/1024x768x0/%:
	./bin/generate.py thumb "$(SOURCE_DIR)" "$(subst $(DEST_DIR)/1024x768x0/,,$@)" 1024x768x0 "$(DEST_DIR)"

$(DEST_DIR)/100x100x1/%:
	./bin/generate.py thumb "$(SOURCE_DIR)" "$(subst $(DEST_DIR)/100x100x1/,,$@)" 100x100x1 "$(DEST_DIR)"

$(DEST_DIR)/118x133x1/%:
	./bin/generate.py thumb "$(SOURCE_DIR)" "$(subst $(DEST_DIR)/118x133x1/,,$@)" 118x133x1 "$(DEST_DIR)"

$(DEST_DIR)/250x250x1/%:
	./bin/generate.py thumb "$(SOURCE_DIR)" "$(subst $(DEST_DIR)/250x250x1/,,$@)" 250x250x1 "$(DEST_DIR)"
