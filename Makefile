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
dest_images := $(subst $(SOURCE_DIR),$(DEST_DIR)/images/1024x768x0,$(images))\
			   $(subst $(SOURCE_DIR),$(DEST_DIR)/images/100x100x1,$(images))\
			   $(subst $(SOURCE_DIR),$(DEST_DIR)/images/512x269x1,$(images))\
			   $(subst $(SOURCE_DIR),$(DEST_DIR)/images/118x133x1,$(images))\
			   $(subst $(SOURCE_DIR),$(DEST_DIR)/images/250x250x1,$(images))
dest_images := $(subst :,\:,$(dest_images))

.PHONY: all
all: $(dest_images)
	./bin/generate.py json '$(JSON_FILE)' $(SOURCE_DIR)
	cp -r www/assets www/styles "$(DEST_DIR)"
	./bin/generator '$(JSON_FILE)' "$(DEST_DIR)"

build:
	go build -o ./bin/generator ./bin/generate_pages.go

$(DEST_DIR)/images/1024x768x0/%:
	./bin/generate.py thumb "$(SOURCE_DIR)/" "$(subst $(DEST_DIR)/images/1024x768x0/,,$@)" 1024x768x0 "$(DEST_DIR)/images/"

$(DEST_DIR)/images/100x100x1/%:
	./bin/generate.py thumb "$(SOURCE_DIR)/" "$(subst $(DEST_DIR)/images/100x100x1/,,$@)" 100x100x1 "$(DEST_DIR)/images/"

$(DEST_DIR)/images/512x269x1/%:
	./bin/generate.py thumb "$(SOURCE_DIR)/" "$(subst $(DEST_DIR)/images/512x269x1/,,$@)" 512x269x1 "$(DEST_DIR)/images/"

$(DEST_DIR)/images/118x133x1/%:
	./bin/generate.py thumb "$(SOURCE_DIR)/" "$(subst $(DEST_DIR)/images/118x133x1/,,$@)" 118x133x1 "$(DEST_DIR)/images/"

$(DEST_DIR)/images/250x250x1/%:
	./bin/generate.py thumb "$(SOURCE_DIR)/" "$(subst $(DEST_DIR)/images/250x250x1/,,$@)" 250x250x1 "$(DEST_DIR)/images/"
