distribution_path = ./distribution
binary_name = ic20
ld_flags = -s -w

all: distribute

generate:
	cd generators && ./cities_generator.py > ../game/cities_generated.go
	cd generators && ./pathogens_generator.py > ../game/pathogens_generated.go

distribute:
	rm -rf $(distribution_path)
	GOOS=windows go build -ldflags="$(ld_flags)" -o $(distribution_path)/$(binary_name)_windows.exe
	GOOS=darwin go build -ldflags="$(ld_flags)" -o $(distribution_path)/$(binary_name)_darwin
	GOOS=linux go build -ldflags="$(ld_flags)" -o $(distribution_path)/$(binary_name)_linux
	cp -r ./example $(distribution_path)
	rm -rf $(distribution_path)/example/__pycache__
	cd $(distribution_path) && zip -r ./$(binary_name)-$(shell date -u "+%s").zip ./$(binary_name)_* ./example
