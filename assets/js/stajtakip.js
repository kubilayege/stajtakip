// https://stackoverflow.com/questions/6982692/html5-input-type-date-default-value-to-today
Date.prototype.toDateInputValue = (function() {
    var local = new Date(this);
    local.setMinutes(this.getMinutes() - this.getTimezoneOffset());
    return local.toJSON().slice(0,10);
});

$(function() {
  $('.datePicker').each(function() {
    if ($(this).val() == "") {
      $(this).val(new Date().toDateInputValue());
    }
  });

  $('#stajToplamGun').keydown(UyariKontrol);
  $('#stajToplamGun').change(UyariKontrol);
  $('#stajSinif').keydown(UyariKontrol);
  $('#stajSinif').change(UyariKontrol);
  $('#stajKonu').change(UyariKontrol);
  $('#stajKurumAdi').change(UyariKontrol);
});

function UyariKontrol() {
  var gun = parseInt($('#stajToplamGun').val());
  var sinif = parseInt($('#stajSinif').val());
  var konu = $('#stajKonu').val().toLowerCase();
  var kurum = $('#stajKurumAdi').val().toLowerCase();

  if (isNaN(gun))
    return;

  if (gun < 15) {
    $("#stajUyari15Gun").show();
  } else {
    $("#stajUyari15Gun").hide();
  }

  if (sinif == 2 && gun > 25) {
    $("#stajUyari25Gun").show();
  } else {
    $("#stajUyari25Gun").hide();
  }

  if (konu != "arge" && gun > 40) {
    $("#stajUyariArgeDegil").show();
  } else {
    $("#stajUyariArgeDegil").hide();
  }

  if (konu == "arge" && gun > 60) {
    $("#stajUyariArge").show();
  } else {
    $("#stajUyariArge").hide();
  }

  for (var i=0; i<stajlar.length; i++) {
    if (stajlar[i].KurumAdi.toLowerCase() == kurum.toLowerCase() && gun+stajlar[i].ToplamGun > 40) {
      $("#stajUyariAyniKurum").show();
      return;
    }
  }

  $("#stajUyariAyniKurum").hide();
}
