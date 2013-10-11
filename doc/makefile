
NAME=paper
OUTPUT=$(NAME).pdf

all: show 

compile:
	pdflatex $(NAME).tex 
#	bibtex $(NAME)
#	pdflatex $(NAME).tex 
#	pdflatex $(NAME).tex 

show: compile
	evince $(OUTPUT) 2>/dev/null &

clean: 
	rm -fr *.blg
	rm -fr *.bbl
	rm -fr *.aux
	rm -fr *.log
