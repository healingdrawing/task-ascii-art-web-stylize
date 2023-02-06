# task-ascii-art-web-stylize
grit:lab Åland Islands 2022

- [Usage: how to run](#usage-how-to-run)  

- [Description](#description)  

- [Authors](#authors)  

- [Implementation details: algorithm](#implementation-details-algorithm)  

## Usage: how to run  

- clone the repository `git clone https://01.gritlab.ax/git/natkim/ascii-art-web.git`;
- terminal: `go run .`;
- open web browser and go to `localhost:8080`;
- type the text you want to convert to the art;
- choose the banner type using radio button switcher;
- press the `convert` button to see the result.
F.e.:  
<img src="readme_files/wnln.png" width=340px>
- if something goes wrong, press the `GO BACK` button.

## Description  

"ascii-art-web" is a web application for converting text consisting of printable standard ASCII characters (including `space`) into a graphic representation of that text.

## Authors  

- [@nattikim](https://github.com/nattikim)  
- [healingdrawing.github.io](https://healingdrawing.github.io)  
- [cenk-idris](https://github.com/cenk-idris)

## Implementation details: algorithm  

The `main.go` file of the root level folder is http server implementation.  

After the execution of the `go run .` command in terminal of the project root level, the server is listening to port `8080`.  

Opening the browser window on the local machine using url `localhost:8080` demonstrates the web application graphical user interface.

After clicking the `convert` button, the filled form data will be sent to the server, where the inputed text will be converted into graphic representation.

The properly formatted graphic representation will be retun back to the client side (into browser window).

In case of some error, a properly formatted error page will be sent to the client side.

In case of unsupported characters are used in the input text, a warning notification will be sent to the client side.  

The `ascii-art` folder includes the package for converting the text into graphic representation.  

The `templates` folder includes the `static` folder with `Raleway-Thin.ttf` and `style.css` files, which are used for stylising of html template files such as `error.html`, `index.html` and `template.html`.  

The `readme_files` folder does not provide any functionality for the `ascii-art-web` application, and is only used to store resources for the `README.MD` file.  

To test server status, run `go test` command in terminal, that will execute `main_test.go` file. Status `200` means that server is ready to manage requests.  

Flag `-debug` can be used to print additional information of server logs. Examples: `go run . -debug` or `go test -debug`.    
