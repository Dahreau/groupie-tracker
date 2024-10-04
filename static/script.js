function hiddenInputs(value){
    var searchInput = document.getElementById("search")
    var searchLocationInput = document.getElementById("searchLocation")
    var creationDateInput = document.getElementById("creationDate")
    var dateInput = document.getElementById("date")
    if (value ==="date" || value ==="firstAlbum"){
        dateInput.classList.remove("hidden")
        searchInput.classList.add("hidden")
        searchLocationInput.classList.add("hidden")
        creationDateInput.classList.add("hidden")
        searchInput.textContent=""
        searchLocationInput.textContent=""
    }else if (value ==="artist"){
        searchInput.classList.remove("hidden")
        dateInput.classList.add("hidden")
        searchLocationInput.classList.add("hidden")
        creationDateInput.classList.add("hidden")
        searchLocationInput.textContent=""
    }else if (value ==="location"){
        searchLocationInput.classList.remove("hidden")
        searchInput.classList.add("hidden")
        dateInput.classList.add("hidden")
        creationDateInput.classList.add("hidden")
        searchInput.innerText=""
    }else if (value ==="creationDate"){
        creationDateInput.classList.remove("hidden")
        searchInput.classList.add("hidden")
        searchLocationInput.classList.add("hidden")
        dateInput.classList.add("hidden")
        searchInput.textContent=""
        searchLocationInput.textContent=""
    }
    
}

    document.getElementById("search").addEventListener('input', function() {
        var query = this.value;
    if (query.length >=1) {
        document.getElementById("search").setAttribute("list", "suggestions");
    }else{
        document.getElementById("search").removeAttribute("list")
    }
})
    document.getElementById("searchLocation").addEventListener('input', function() {
        var query = this.value;
    if (query.length >1) {
        document.getElementById("searchLocation").setAttribute("list", "suggestionsLocations");
    }else{
        document.getElementById("searchLocation").removeAttribute("list")
    }
    });
    document.getElementById("creationDate").addEventListener('input', function() {
        var query = this.value;
    if (query.length >1) {
        document.getElementById("creationDate").setAttribute("list", "suggestionsCreationDate");
    }else{
        document.getElementById("creationDate").removeAttribute("list")
    }
});

const minDate = document.getElementById('minDate');
const outputMinDate = document.getElementById("outputMinDate");
const maxDate = document.getElementById('maxDate');
const outputMaxDate = document.getElementById("outputMaxDate");

const minDateFirstAlbum = document.getElementById('minDateFirstAlbum');
const outputMinDateFirstAlbum = document.getElementById("outputMinDateFirstAlbum");
const maxDateFirstAlbum = document.getElementById('maxDateFirstAlbum');
const outputMaxDateFirstAlbum = document.getElementById("outputMaxDateFirstAlbum");



function updateMaxSlider(newValue) {
  if (newValue > maxDate.value) {
    maxDate.value = newValue;
    outputMaxDate.innerText = newValue
}
}

function updateMinSlider(newValue) {
  if (newValue < minDate.value) {
    minDate.value = newValue;
    outputMinDate.innerText = newValue;
}
}

function updateMaxSliderFirstAlbum(newValue) {
    if (newValue > maxDateFirstAlbum.value) {
      maxDateFirstAlbum.value = newValue;
      outputMaxDateFirstAlbum.innerText = newValue
  }
  }
  
  function updateMinSliderFirstAlbum(newValue) {
    if (newValue < minDateFirstAlbum.value) {
      minDateFirstAlbum.value = newValue;
      outputMinDateFirstAlbum.innerText = newValue;
  }
  }
minDate.addEventListener('input', () => {
  outputMinDate.innerText = minDate.value
  updateMaxSlider(minDate.value);
});

maxDate.addEventListener('input', () => {
  outputMaxDate.innerText = maxDate.value;
  updateMinSlider(maxDate.value);
});


minDateFirstAlbum.addEventListener('input', () => {
    outputMinDateFirstAlbum.innerText = minDateFirstAlbum.value
    updateMaxSliderFirstAlbum(minDateFirstAlbum.value);
  });
  
  maxDateFirstAlbum.addEventListener('input', () => {
    outputMaxDateFirstAlbum.innerText = maxDateFirstAlbum.value;
    updateMinSliderFirstAlbum(maxDateFirstAlbum.value);
});