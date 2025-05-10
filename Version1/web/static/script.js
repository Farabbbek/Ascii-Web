document.addEventListener("DOMContentLoaded", () => {
  console.log("%c DOMContentLoaded event fired", "color: green; font-weight: bold");

  if (window.location.hash) {
    console.log("URL contains fragment:", window.location.hash);
    window.location.href = "/nonexistent";
    return;
  }

  const divider = document.getElementById("divider");
  const sidebar = document.getElementById("sidebar");
  const container = document.querySelector(".container");

  if (divider && sidebar && container) {
    let isDragging = false;

    divider.addEventListener("mousedown", (e) => {
      isDragging = true;
      document.body.style.cursor = "ew-resize";
    });

    document.addEventListener("mousemove", (e) => {
      if (!isDragging) return;

      const containerRect = container.getBoundingClientRect();
      const newSidebarWidth = containerRect.right - e.clientX;

      const minWidth = 200;
      const maxWidth = containerRect.width * 0.6;

      if (newSidebarWidth >= minWidth && newSidebarWidth <= maxWidth) {
        sidebar.style.flex = `0 0 ${newSidebarWidth}px`;
      }
    });

    document.addEventListener("mouseup", () => {
      isDragging = false;
      document.body.style.cursor = "default";
    });

    divider.addEventListener("dragstart", (e) => e.preventDefault());
  } else {
    console.error("Divider elements not found:", {
      divider: !!divider,
      sidebar: !!sidebar,
      container: !!container,
    });
  }

  const form = document.getElementById("ascii-art-form");
  const outputElement = document.getElementById("ascii-art-output");
  const textInput = document.getElementById("text-input");
  const bannerSelect = document.getElementById("banner-select");
  const generateBtn = document.getElementById("generate-btn");
  const viewFullPageBtn = document.getElementById("view-full-page");
  const scrollHint = document.getElementById("scroll-hint");
  const copyBtn = document.getElementById("copy-btn");
  const asciiContainer = document.getElementById("ascii-art-container");

  if (!form || !outputElement || !textInput || !bannerSelect || !generateBtn || !viewFullPageBtn) {
    console.error("One or more form elements not found:", {
      form: !!form,
      outputElement: !!outputElement,
      textInput: !!textInput,
      bannerSelect: !!bannerSelect,
      generateBtn: !!generateBtn,
      viewFullPageBtn: !!viewFullPageBtn,
    });
    return;
  }

  console.log("%c Attaching submit event listener to form", "color: blue; font-weight: bold");

  copyBtn.addEventListener("click", () => {
    const textToCopy = outputElement.value;
    navigator.clipboard
      .writeText(textToCopy)
      .then(() => {
        copyBtn.textContent = "Скопировано!";
        setTimeout(() => {
          copyBtn.textContent = "Копировать ASCII-арт";
        }, 2000);
      })
      .catch((err) => {
        console.error("Ошибка при копировании: ", err);
        copyBtn.textContent = "Ошибка копирования";
        setTimeout(() => {
          copyBtn.textContent = "Копировать ASCII-арт";
        }, 2000);
      });
  });

  form.addEventListener("submit", (e) => {
    e.preventDefault();
    console.log("%c Form submit triggered", "color: purple; font-weight: bold");

    generateBtn.disabled = true;
    generateBtn.textContent = "Generating...";

    let text = textInput.value;
    const banner = bannerSelect.value;

    console.log("Raw input text:", JSON.stringify(text));
    text = text.replace(/\\n/g, "\n");
    console.log("Processed input text:", JSON.stringify(text));
    console.log("Form data being sent:", { text, banner });

    if (!text) {
      console.log("Validation failed: Text is empty");
      window.location.href = "/400";
      generateBtn.disabled = false;
      generateBtn.textContent = "Generate ASCII Art";
      return;
    }

    if (!/^[ -~\n]*$/.test(text)) {
      console.log("Validation failed: Text contains invalid characters");
      window.location.href = "/400";
      generateBtn.disabled = false;
      generateBtn.textContent = "Generate ASCII Art";
      return;
    }

    const validBanners = ["standard", "shadow", "thinkertoy"];
    if (!validBanners.includes(banner)) {
      console.log("Validation failed: Invalid banner choice");
      window.location.href = "/400";
      generateBtn.disabled = false;
      generateBtn.textContent = "Generate ASCII Art";
      return;
    }

    const params = new URLSearchParams();
    params.append("text", text);
    params.append("banner", banner);

    console.log("URLSearchParams:", params.toString());

    fetch("/generate-ascii", {
      method: "POST",
      body: params,
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      redirect: "manual",
    })
      .then((response) => {
        console.log("Response status:", response.status, "Response type:", response.type);
        if (response.type === "opaqueredirect") {
          console.log("Redirecting to /400");
          window.location.href = "/400";
          return;
        }
        if (!response.ok) {
          if (response.status === 500) {
            throw new Error("Внутренняя ошибка сервера. Пожалуйста, попробуйте позже.");
          }
          return response.text().then((text) => {
            throw new Error(`Ошибка сервера: ${response.status} - ${text || "Нет деталей"}`);
          });
        }
        return response.text();
      })
      .then((data) => {
        console.log("Received data:", JSON.stringify(data));
        const normalizedData = data.replace(/\r\n|\r/g, "\n");

        outputElement.value = normalizedData;
        outputElement.style.whiteSpace = "pre";
        outputElement.style.overflowX = "scroll";
        outputElement.style.fontFamily = "'Courier New', Courier, monospace";
        outputElement.style.display = "block";
        outputElement.style.color = ""; // Сбрасываем цвет
        outputElement.style.fontWeight = ""; // Сбрасываем жирность

        const tempDiv = document.createElement("div");
        tempDiv.style.fontFamily = "'Courier New', Courier, monospace";
        tempDiv.style.whiteSpace = "pre";
        tempDiv.style.position = "absolute";
        tempDiv.style.visibility = "hidden";
        tempDiv.style.width = "auto";

        const lines = normalizedData.split("\n");
        const longestLine = lines.reduce((a, b) => (a.length > b.length ? a : b), "");
        tempDiv.textContent = longestLine;
        document.body.appendChild(tempDiv);

        const contentWidth = tempDiv.offsetWidth;
        const containerWidth = asciiContainer.offsetWidth;
        document.body.removeChild(tempDiv);

        if (contentWidth > containerWidth - 30) {
          scrollHint.style.display = "block";
          setTimeout(() => {
            outputElement.style.display = "none";
            setTimeout(() => {
              outputElement.style.display = "block";
            }, 10);
          }, 10);
        } else {
          scrollHint.style.display = "none";
        }

        const lineCount = lines.length;
        const minHeight = Math.min(Math.max(lineCount * 20, 100), 600);
        outputElement.style.minHeight = `${minHeight}px`;

        viewFullPageBtn.style.display = "block";
        copyBtn.style.display = "block";

        asciiContainer.scrollIntoView({ behavior: "smooth" });
      })
      .catch((error) => {
        console.error("Error generating ASCII art:", error);
        outputElement.value = error.message;
        outputElement.style.display = "block";
        outputElement.style.color = "red"; // Красный текст для ошибки
        outputElement.style.fontWeight = "bold"; // Жирный текст для ошибки
        scrollHint.style.display = "none";
        copyBtn.style.display = "none";
        asciiContainer.scrollIntoView({ behavior: "smooth" });
      })
      .finally(() => {
        console.log("Request completed");
        generateBtn.disabled = false;
        generateBtn.textContent = "Generate ASCII Art";
      });
  });

  textInput.addEventListener("input", () => {
    console.log("Text input changed, clearing output");
    outputElement.style.display = "none";
    outputElement.value = "";
    scrollHint.style.display = "none";
    copyBtn.style.display = "none";
    viewFullPageBtn.style.display = "none";
  });

  generateBtn.addEventListener("click", (e) => {
    e.preventDefault();
    console.log("%c Generate button clicked", "color: orange; font-weight: bold");
    form.dispatchEvent(new Event("submit", { cancelable: true }));
  });

  viewFullPageBtn.addEventListener("click", () => {
    const asciiArt = outputElement.value;
    const newWindow = window.open("", "_blank");
    newWindow.document.write(`
      <!DOCTYPE html>
      <html lang="en">
      <head>
          <meta charset="UTF-8">
          <meta name="viewport" content="width=device-width, initial-scale=1.0">
          <title>ASCII Art Result</title>
          <style>
              body {
                  margin: 0;
                  padding: 20px;
                  background: #1a191c;
                  color: #ffffff;
                  font-family: sans-serif;
                  height: 100vh;
                  display: flex;
                  flex-direction: column;
              }
              h1 {
                  font-family: monospace;
                  font-size: 2rem;
                  margin-bottom: 20px;
              }
              .ascii-container {
                  width: 100%;
                  height: calc(100% - 120px);
                  overflow: auto;
                  background: #272529;
                  border: 1px solid #bdbdbd;
                  padding: 0;
                  box-sizing: border-box;
                  position: relative;
              }
              pre {
                  margin: 0;
                  padding: 15px;
                  white-space: pre;
                  font-family: 'Courier New', Courier, monospace;
                  font-size: 1rem;
                  color: #ffffff;
                  line-height: 1.2;
                  tab-size: 4;
              }
              .controls {
                  display: flex;
                  justify-content: space-between;
                  margin-top: 10px;
                  flex-wrap: wrap;
                  gap: 10px;
              }
              a, button {
                  color: #1e90ff;
                  text-decoration: none;
                  font-family: monospace;
                  background: none;
                  border: none;
                  cursor: pointer;
                  padding: 5px 10px;
              }
              a:hover, button:hover {
                  text-decoration: underline;
              }
              .zoom-controls {
                  display: flex;
                  gap: 10px;
                  align-items: center;
              }
              .zoom-level {
                  color: #ffffff;
                  font-family: monospace;
              }
              .copy-btn {
                  background-color: #00D4A1;
                  color: #ffffff;
                  border: none;
                  border-radius: 4px;
                  padding: 5px 15px;
              }
              .copy-btn:hover {
                  background-color: #00b389;
                  text-decoration: none;
              }
              .info-bar {
                  display: flex;
                  justify-content: space-between;
                  font-family: monospace;
                  font-size: 0.8rem;
                  color: #bdbdbd;
                  margin-top: 5px;
              }
          </style>
      </head>
      <body>
          <h1>ASCII Art Result</h1>
          <div class="ascii-container">
              <pre id="ascii-output">${asciiArt}</pre>
          </div>
          <div class="info-bar">
              <span id="size-info">Используйте колесико мыши для прокрутки</span>
              <span id="char-count">${asciiArt.length} символов</span>
          </div>
          <div class="controls">
              <a href="/">Вернуться на главную</a>
              <div class="zoom-controls">
                  <button id="zoom-out">Уменьшить</button>
                  <span class="zoom-level" id="zoom-level">100%</span>
                  <button id="zoom-in">Увеличить</button>
              </div>
              <button class="copy-btn" id="copy-btn">Копировать</button>
          </div>
          <script>
              const output = document.getElementById('ascii-output');
              const zoomLevel = document.getElementById('zoom-level');
              const zoomIn = document.getElementById('zoom-in');
              const zoomOut = document.getElementById('zoom-out');
              const copyBtn = document.getElementById('copy-btn');
              const sizeInfo = document.getElementById('size-info');
              
              let currentZoom = 100;
              
              const updateSizeInfo = () => {
                  const lines = output.textContent.split('\\n');
                  const maxLineLength = lines.reduce((max, line) => 
                      Math.max(max, line.length), 0);
                  sizeInfo.textContent = 
                      \`\${lines.length} строк, макс. ширина: \${maxLineLength} символов\`;
              };
              
              updateSizeInfo();
              
              zoomIn.addEventListener('click', () => {
                  if (currentZoom < 200) {
                      currentZoom += 10;
                      output.style.fontSize = (currentZoom / 100) + 'rem';
                      zoomLevel.textContent = currentZoom + '%';
                  }
              });
              
              zoomOut.addEventListener('click', () => {
                  if (currentZoom > 50) {
                      currentZoom -= 10;
                      output.style.fontSize = (currentZoom / 100) + 'rem';
                      zoomLevel.textContent = currentZoom + '%';
                  }
              });
              
              copyBtn.addEventListener('click', () => {
                  const textToCopy = output.textContent;
                  navigator.clipboard.writeText(textToCopy)
                      .then(() => {
                          copyBtn.textContent = 'Скопировано!';
                          setTimeout(() => {
                              copyBtn.textContent = 'Копировать';
                          }, 2000);
                      })
                      .catch(err => {
                          console.error('Ошибка при копировании: ', err);
                          copyBtn.textContent = 'Ошибка';
                          setTimeout(() => {
                              copyBtn.textContent = 'Копировать';
                          }, 2000);
                      });
              });
              
              document.addEventListener('keydown', (e) => {
                  if (e.ctrlKey) {
                      if (e.key === '+' || e.key === '=') {
                          e.preventDefault();
                          zoomIn.click();
                      } else if (e.key === '-') {
                          e.preventDefault();
                          zoomOut.click();
                      }
                  }
              });
          </script>
      </body>
      </html>
    `);
    newWindow.document.close();
  });

  const goBackButton = document.querySelector(".animate-01.inlineBlock-01.fMono-01.medium-01.regular-01.fs6-01");
  if (goBackButton) {
    goBackButton.addEventListener("click", () => {
      window.history.back();
    });
  }
});