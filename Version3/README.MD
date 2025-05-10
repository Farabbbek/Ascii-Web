ASCII Art Web

📌 Project Overview

This project is a web-based ASCII Art Generator that allows users to input text and convert it into stylish ASCII representations. The interface provides various formatting options and a user-friendly experience.

🚀 Features

Text Input: Users can type any text to generate ASCII art.

Style Selection: Multiple ASCII font styles.

Live Preview: Output is displayed in real-time.

Responsive Design: Works well on different screen sizes.

Custom Styling: Clean UI with CSS enhancements.

📂 Project Structure
    ascii-art-web/
    ├── main/
    │   ├── main.go
    │   ├── ascii_art.go
    │   └── index.html
    │   └── style.css
    └── banner/
        ├── standard.txt
        ├── shadow.txt
        └── thinkertoy.txt


🛠️ Technologies Used

HTML – Structure

CSS – Styling

GO – Logic for ASCII conversion


🔧 Setup Instructions

Navigate to the project folder:

cd ascii-art-web/ascii-web/main

Run the application:
 
go run main.go ascii_art.go

Open your web browser and navigate to:
    http://localhost:8080



📜 Code Explanation

main.go

This file contains the main server code. It sets up the HTTP server, handles requests, and serves static files.

ascii_art.go

This file contains the logic for loading font files and generating ASCII art from text input.

index.html

This file contains the frontend code. It includes the HTML structure, CSS for styling, and JavaScript for handling user input and making requests to the server.

style.css

This file contains the CSS styles for the frontend, ensuring a clean and responsive design.


🎨 Usage Guide

Enter text into the input box.

Choose a desired ASCII style.

View the generated ASCII output.

Copy and use it anywhere!


👤 Author

Developed by Tusup AL-Farabi

