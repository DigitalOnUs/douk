function UploadFile(){
    let input = document.querySelector('#file-upload input[type=file]');
    if (input.files.length > 0){
        const fileName = document.querySelector('#file-upload .file-name');
        fileName .textContent = input.files[0].name;

        // set content for zone 
        const display = document.getElementById("input-doc");
        console.log(input.files[0].name)
        console.log("here we are");
        
        const reader = new FileReader();
        reader.onload = function fileReadCompleted() {
          // when the reader is done, the content is in reader.result.
          console.log(reader.result);
          display.value =reader.result
        };
        reader.readAsText(input.files[0]);
    }
}

function SendToBackend(){
    console.log("there you are")
    // validating type file 
    const fileName = document.querySelector('#file-upload .file-name')
    const extension = fileName.textContent.match(/\.[0-9a-z]+$/i);
    // looks awful
    if (extension){ 
        (extension[0] === '.json' || extension[0] === '.hcl') ? console.log(extension): console.log("nel ...");
        return
    }
    window.alert("Supported file extensions are json/hcl");
}