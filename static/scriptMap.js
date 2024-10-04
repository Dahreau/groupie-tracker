
function initMap(){
    const apiKey = 'R4ekZroMZ5pJPn1tBISF'; // Replace with your MapTiler API key
    maptilersdk.config.apiKey = apiKey;
    const map = new maptilersdk.Map({
        container: 'map', // container's id or the HTML element in which the SDK will render the map
        style: maptilersdk.MapStyle.STREETS, // starting position [lng, lat]
        zoom: 9 // starting zoom
    });
    map.addControl(new maptilersdk.NavigationControl()); // add navigation controls
    map.addControl(new maptilersdk.GeolocateControl());
    const coords = extractCoordinates();
    if (coords.length > 0) {
        const bounds = new maptilersdk.LngLatBounds();

        coords.forEach(coord => {
            addMarker(map, coord.Longitude, coord.Latitude,coord.Location);
            bounds.extend([coord.Longitude, coord.Latitude]);
        });

        map.fitBounds(bounds, {
            padding: 50
        });
    }
    
    return map;
}
function addMarker(map, lng, lat, text) {
    text = text.replaceAll("_"," ");
    text = text.replaceAll("-"," - ")
    text = capitalizeWords(text);
    const popup = new maptilersdk.Popup({ offset: 25 }).setText(text);
    new maptilersdk.Marker()
        .setLngLat([lng, lat])
        .setPopup(popup)
        .addTo(map);
}
    // Initialiser la carte quand le DOM est chargé
    document.addEventListener('DOMContentLoaded', () => {
        const map = initMap();
    });

    function extractCoordinates() {
        const longitudes = document.querySelectorAll('#longitude li');
        const latitudes = document.querySelectorAll('#latitude li');
        const locations = document.querySelectorAll('#location li');

        const coordinates = [];

        for (let i = 0; i < longitudes.length; i++) {
            const longitude = parseFloat(longitudes[i].textContent.trim());
            const latitude = parseFloat(latitudes[i].textContent.trim());
            const location = locations[i].textContent.trim();

            coordinates.push({
                Location: location,
                Longitude: longitude,
                Latitude: latitude
            });
        }

        return coordinates;
    }

    function capitalizeWords(str) {
        return str.split(' ').map(word => {
            return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
        }).join(' ');
    }