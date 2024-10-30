window.addEventListener("load", function () {
  let pictureLabel = document.getElementById("picture-label");
  let pictureInput = document.getElementById("picture-input");
  let uploadButton = document.getElementById("upload-button");
  let undo = document.getElementById("undo");

  pictureLabel.addEventListener("click", () => {
    pictureInput.click();
  });

  pictureInput.addEventListener("change", () => {
    pictureLabel.style.display = "none";
    uploadButton.style.display = "flex";
  });

  undo.addEventListener("click", () => {
    pictureLabel.style.display = "inline";
    uploadButton.style.display = "none";
  });
});
