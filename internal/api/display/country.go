package display

import (
	"encoding/json"
)

const conutryJson = ` [
    {
        "name": "Afghanistan",
        "abr": "AF",
        "img": "https://cdn.countryflags.com/thumbs/afghanistan/flag-square-250.png"
    },
    {
        "name": "Åland Islands",
        "abr": "AX",
        "img": "https://cdn.countryflags.com/thumbs/aland-Islands/flag-square-250.png"
    },
    {
        "name": "Albania",
        "abr": "AL",
        "img": "https://cdn.countryflags.com/thumbs/albania/flag-square-250.png"
    },
    {
        "name": "Algeria",
        "abr": "DZ",
        "img": "https://cdn.countryflags.com/thumbs/algeria/flag-square-250.png"
    },
    {
        "name": "American Samoa",
        "abr": "AS",
        "img": ""
    },
    {
        "name": "Andorra",
        "abr": "AD",
        "img": "https://cdn.countryflags.com/thumbs/andorra/flag-square-250.png"
    },
    {
        "name": "Angola",
        "abr": "AO",
        "img": "https://cdn.countryflags.com/thumbs/angola/flag-square-250.png"
    },
    {
        "name": "Anguilla",
        "abr": "AI",
        "img": ""
    },
    {
        "name": "Antarctica",
        "abr": "AQ",
        "img": ""
    },
    {
        "name": "Antigua and Barbuda",
        "abr": "AG",
        "img": "https://cdn.countryflags.com/thumbs/antigua-and-barbuda/flag-square-250.png"
    },
    {
        "name": "Argentina",
        "abr": "AR",
        "img": "https://cdn.countryflags.com/thumbs/argentina/flag-square-250.png"
    },
    {
        "name": "Armenia",
        "abr": "AM",
        "img": "https://cdn.countryflags.com/thumbs/armenia/flag-square-250.png"
    },
    {
        "name": "Aruba",
        "abr": "AW",
        "img": "https://cdn.countryflags.com/thumbs/aruba/flag-square-250.png"
    },
    {
        "name": "Australia",
        "abr": "AU",
        "img": "https://cdn.countryflags.com/thumbs/australia/flag-square-250.png"
    },
    {
        "name": "Austria",
        "abr": "AT",
        "img": "https://cdn.countryflags.com/thumbs/austria/flag-square-250.png"
    },
    {
        "name": "Azerbaijan",
        "abr": "AZ",
        "img": "https://cdn.countryflags.com/thumbs/azerbaijan/flag-square-250.png"
    },
    {
        "name": "Bahamas",
        "abr": "BS",
        "img": "https://cdn.countryflags.com/thumbs/bahamas/flag-square-250.png"
    },
    {
        "name": "Bahrain",
        "abr": "BH",
        "img": "https://cdn.countryflags.com/thumbs/bahrain/flag-square-250.png"
    },
    {
        "name": "Bangladesh",
        "abr": "BD",
        "img": "https://cdn.countryflags.com/thumbs/bangladesh/flag-square-250.png"
    },
    {
        "name": "Barbados",
        "abr": "BB",
        "img": "https://cdn.countryflags.com/thumbs/barbados/flag-square-250.png"
    },
    {
        "name": "Belarus",
        "abr": "BY",
        "img": "https://cdn.countryflags.com/thumbs/belarus/flag-square-250.png"
    },
    {
        "name": "Belgium",
        "abr": "BE",
        "img": "https://cdn.countryflags.com/thumbs/belgium/flag-square-250.png"
    },
    {
        "name": "Belize",
        "abr": "BZ",
        "img": "https://cdn.countryflags.com/thumbs/belize/flag-square-250.png"
    },
    {
        "name": "Benin",
        "abr": "BJ",
        "img": "https://cdn.countryflags.com/thumbs/benin/flag-square-250.png"
    },
    {
        "name": "Bermuda",
        "abr": "BM",
        "img": ""
    },
    {
        "name": "Bhutan",
        "abr": "BT",
        "img": "https://cdn.countryflags.com/thumbs/bhutan/flag-square-250.png"
    },
    {
        "name": "Bolivia (Plurinational State of)",
        "abr": "BO",
        "img": "https://cdn.countryflags.com/thumbs/bolivia/flag-square-250.png"
    },
    {
        "name": "Bonaire, Sint Eustatius and Saba",
        "abr": "BQ",
        "img": ""
    },
    {
        "name": "Bosnia and Herzegovina",
        "abr": "BA",
        "img": "https://cdn.countryflags.com/thumbs/bosnia-and-herzegovina/flag-square-250.png"
    },
    {
        "name": "Botswana",
        "abr": "BW",
        "img": "https://cdn.countryflags.com/thumbs/botswana/flag-square-250.png"
    },
    {
        "name": "Bouvet Island",
        "abr": "BV",
        "img": ""
    },
    {
        "name": "Brazil",
        "abr": "BR",
        "img": "https://cdn.countryflags.com/thumbs/brazil/flag-square-250.png"
    },
    {
        "name": "British Indian Ocean Territory",
        "abr": "IO",
        "img": ""
    },
    {
        "name": "Brunei Darussalam",
        "abr": "BN",
        "img": ""
    },
    {
        "name": "Bulgaria",
        "abr": "BG",
        "img": "https://cdn.countryflags.com/thumbs/bulgaria/flag-square-250.png"
    },
    {
        "name": "Burkina Faso",
        "abr": "BF",
        "img": "https://cdn.countryflags.com/thumbs/burkina-faso/flag-square-250.png"
    },
    {
        "name": "Burundi",
        "abr": "BI",
        "img": "https://cdn.countryflags.com/thumbs/burundi/flag-square-250.png"
    },
    {
        "name": "Cabo Verde",
        "abr": "CV",
        "img": ""
    },
    {
        "name": "Cambodia",
        "abr": "KH",
        "img": "https://cdn.countryflags.com/thumbs/cambodia/flag-square-250.png"
    },
    {
        "name": "Cameroon",
        "abr": "CM",
        "img": "https://cdn.countryflags.com/thumbs/cameroon/flag-square-250.png"
    },
    {
        "name": "Canada",
        "abr": "CA",
        "img": "https://cdn.countryflags.com/thumbs/canada/flag-square-250.png"
    },
    {
        "name": "Cayman Islands",
        "abr": "KY",
        "img": ""
    },
    {
        "name": "Central African Republic",
        "abr": "CF",
        "img": ""
    },
    {
        "name": "Chad",
        "abr": "TD",
        "img": "https://cdn.countryflags.com/thumbs/chad/flag-square-250.png"
    },
    {
        "name": "Chile",
        "abr": "CL",
        "img": "https://cdn.countryflags.com/thumbs/chile/flag-square-250.png"
    },
    {
        "name": "China",
        "abr": "CN",
        "img": "https://cdn.countryflags.com/thumbs/china/flag-square-250.png"
    },
    {
        "name": "Christmas Island",
        "abr": "CX",
        "img": ""
    },
    {
        "name": "Cocos (Keeling) Islands",
        "abr": "CC",
        "img": ""
    },
    {
        "name": "Colombia",
        "abr": "CO",
        "img": "https://cdn.countryflags.com/thumbs/colombia/flag-square-250.png"
    },
    {
        "name": "Comoros",
        "abr": "KM",
        "img": "https://cdn.countryflags.com/thumbs/comoros/flag-square-250.png"
    },
    {
        "name": "Congo",
        "abr": "CG",
        "img": ""
    },
    {
        "name": "Congo (Democratic Republic of the)",
        "abr": "CD",
        "img": ""
    },
    {
        "name": "Cook Islands",
        "abr": "CK",
        "img": ""
    },
    {
        "name": "Costa Rica",
        "abr": "CR",
        "img": "https://cdn.countryflags.com/thumbs/costa-rica/flag-square-250.png"
    },
    {
        "name": "Côte d'Ivoire",
        "abr": "CI",
        "img": ""
    },
    {
        "name": "Croatia",
        "abr": "HR",
        "img": "https://cdn.countryflags.com/thumbs/croatia/flag-square-250.png"
    },
    {
        "name": "Cuba",
        "abr": "CU",
        "img": "https://cdn.countryflags.com/thumbs/cuba/flag-square-250.png"
    },
    {
        "name": "Curaçao",
        "abr": "CW",
        "img": "https://cdn.countryflags.com/thumbs/curacao/flag-square-250.png"
    },
    {
        "name": "Cyprus",
        "abr": "CY",
        "img": "https://cdn.countryflags.com/thumbs/cyprus/flag-square-250.png"
    },
    {
        "name": "Czechia",
        "abr": "CZ",
        "img": ""
    },
    {
        "name": "Denmark",
        "abr": "DK",
        "img": "https://cdn.countryflags.com/thumbs/denmark/flag-square-250.png"
    },
    {
        "name": "Djibouti",
        "abr": "DJ",
        "img": "https://cdn.countryflags.com/thumbs/djibouti/flag-square-250.png"
    },
    {
        "name": "Dominica",
        "abr": "DM",
        "img": "https://cdn.countryflags.com/thumbs/dominica/flag-square-250.png"
    },
    {
        "name": "Dominican Republic",
        "abr": "DO",
        "img": "https://cdn.countryflags.com/thumbs/dominican-republic/flag-square-250.png"
    },
    {
        "name": "Ecuador",
        "abr": "EC",
        "img": "https://cdn.countryflags.com/thumbs/ecuador/flag-square-250.png"
    },
    {
        "name": "Egypt",
        "abr": "EG",
        "img": "https://cdn.countryflags.com/thumbs/egypt/flag-square-250.png"
    },
    {
        "name": "El Salvador",
        "abr": "SV",
        "img": "https://cdn.countryflags.com/thumbs/el-salvador/flag-square-250.png"
    },
    {
        "name": "Equatorial Guinea",
        "abr": "GQ",
        "img": "https://cdn.countryflags.com/thumbs/equatorial-guinea/flag-square-250.png"
    },
    {
        "name": "Eritrea",
        "abr": "ER",
        "img": "https://cdn.countryflags.com/thumbs/eritrea/flag-square-250.png"
    },
    {
        "name": "Estonia",
        "abr": "EE",
        "img": "https://cdn.countryflags.com/thumbs/estonia/flag-square-250.png"
    },
    {
        "name": "Eswatini",
        "abr": "SZ",
        "img": ""
    },
    {
        "name": "Ethiopia",
        "abr": "ET",
        "img": "https://cdn.countryflags.com/thumbs/ethiopia/flag-square-250.png"
    },
    {
        "name": "Falkland Islands (Malvinas)",
        "abr": "FK",
        "img": ""
    },
    {
        "name": "Faroe Islands",
        "abr": "FO",
        "img": ""
    },
    {
        "name": "Fiji",
        "abr": "FJ",
        "img": "https://cdn.countryflags.com/thumbs/fiji/flag-square-250.png"
    },
    {
        "name": "Finland",
        "abr": "FI",
        "img": "https://cdn.countryflags.com/thumbs/finland/flag-square-250.png"
    },
    {
        "name": "France",
        "abr": "FR",
        "img": "https://cdn.countryflags.com/thumbs/france/flag-square-250.png"
    },
    {
        "name": "French Guiana",
        "abr": "GF",
        "img": ""
    },
    {
        "name": "French Polynesia",
        "abr": "PF",
        "img": ""
    },
    {
        "name": "French Southern Territories",
        "abr": "TF",
        "img": ""
    },
    {
        "name": "Gabon",
        "abr": "GA",
        "img": "https://cdn.countryflags.com/thumbs/gabon/flag-square-250.png"
    },
    {
        "name": "Gambia",
        "abr": "GM",
        "img": ""
    },
    {
        "name": "Georgia",
        "abr": "GE",
        "img": "https://cdn.countryflags.com/thumbs/georgia/flag-square-250.png"
    },
    {
        "name": "Germany",
        "abr": "DE",
        "img": "https://cdn.countryflags.com/thumbs/germany/flag-square-250.png"
    },
    {
        "name": "Ghana",
        "abr": "GH",
        "img": "https://cdn.countryflags.com/thumbs/ghana/flag-square-250.png"
    },
    {
        "name": "Gibraltar",
        "abr": "GI",
        "img": ""
    },
    {
        "name": "Greece",
        "abr": "GR",
        "img": "https://cdn.countryflags.com/thumbs/greece/flag-square-250.png"
    },
    {
        "name": "Greenland",
        "abr": "GL",
        "img": "https://cdn.countryflags.com/thumbs/greenland/flag-square-250.png"
    },
    {
        "name": "Grenada",
        "abr": "GD",
        "img": "https://cdn.countryflags.com/thumbs/grenada/flag-square-250.png"
    },
    {
        "name": "Guadeloupe",
        "abr": "GP",
        "img": "https://cdn.countryflags.com/thumbs/guadeloupe/flag-square-250.png"
    },
    {
        "name": "Guam",
        "abr": "GU",
        "img": "https://cdn.countryflags.com/thumbs/guam/flag-square-250.png"
    },
    {
        "name": "Guatemala",
        "abr": "GT",
        "img": "https://cdn.countryflags.com/thumbs/guatemala/flag-square-250.png"
    },
    {
        "name": "Guernsey",
        "abr": "GG",
        "img": ""
    },
    {
        "name": "Guinea",
        "abr": "GN",
        "img": "https://cdn.countryflags.com/thumbs/guinea/flag-square-250.png"
    },
    {
        "name": "Guinea-Bissau",
        "abr": "GW",
        "img": "https://cdn.countryflags.com/thumbs/guinea-bissau/flag-square-250.png"
    },
    {
        "name": "Guyana",
        "abr": "GY",
        "img": "https://cdn.countryflags.com/thumbs/guyana/flag-square-250.png"
    },
    {
        "name": "Haiti",
        "abr": "HT",
        "img": "https://cdn.countryflags.com/thumbs/haiti/flag-square-250.png"
    },
    {
        "name": "Heard Island and McDonald Islands",
        "abr": "HM",
        "img": ""
    },
    {
        "name": "Holy See",
        "abr": "VA",
        "img": ""
    },
    {
        "name": "Honduras",
        "abr": "HN",
        "img": "https://cdn.countryflags.com/thumbs/honduras/flag-square-250.png"
    },
    {
        "name": "Hong Kong",
        "abr": "HK",
        "img": "https://cdn.countryflags.com/thumbs/hongkong/flag-square-250.png"
    },
    {
        "name": "Hungary",
        "abr": "HU",
        "img": "https://cdn.countryflags.com/thumbs/hungary/flag-square-250.png"
    },
    {
        "name": "Iceland",
        "abr": "IS",
        "img": "https://cdn.countryflags.com/thumbs/iceland/flag-square-250.png"
    },
    {
        "name": "India",
        "abr": "IN",
        "img": "https://cdn.countryflags.com/thumbs/india/flag-square-250.png"
    },
    {
        "name": "Indonesia",
        "abr": "ID",
        "img": "https://cdn.countryflags.com/thumbs/indonesia/flag-square-250.png"
    },
    {
        "name": "Iran (Islamic Republic of)",
        "abr": "IR",
        "img": ""
    },
    {
        "name": "Iraq",
        "abr": "IQ",
        "img": "https://cdn.countryflags.com/thumbs/iraq/flag-square-250.png"
    },
    {
        "name": "Ireland",
        "abr": "IE",
        "img": "https://cdn.countryflags.com/thumbs/ireland/flag-square-250.png"
    },
    {
        "name": "Isle of Man",
        "abr": "IM",
        "img": ""
    },
    {
        "name": "Israel",
        "abr": "IL",
        "img": "https://cdn.countryflags.com/thumbs/israel/flag-square-250.png"
    },
    {
        "name": "Italy",
        "abr": "IT",
        "img": "https://cdn.countryflags.com/thumbs/italy/flag-square-250.png"
    },
    {
        "name": "Jamaica",
        "abr": "JM",
        "img": "https://cdn.countryflags.com/thumbs/jamaica/flag-square-250.png"
    },
    {
        "name": "Japan",
        "abr": "JP",
        "img": "https://cdn.countryflags.com/thumbs/japan/flag-square-250.png"
    },
    {
        "name": "Jersey",
        "abr": "JE",
        "img": ""
    },
    {
        "name": "Jordan",
        "abr": "JO",
        "img": "https://cdn.countryflags.com/thumbs/jordan/flag-square-250.png"
    },
    {
        "name": "Kazakhstan",
        "abr": "KZ",
        "img": "https://cdn.countryflags.com/thumbs/kazakhstan/flag-square-250.png"
    },
    {
        "name": "Kenya",
        "abr": "KE",
        "img": "https://cdn.countryflags.com/thumbs/kenya/flag-square-250.png"
    },
    {
        "name": "Kiribati",
        "abr": "KI",
        "img": "https://cdn.countryflags.com/thumbs/kiribati/flag-square-250.png"
    },
    {
        "name": "Korea (Democratic People's Republic of)",
        "abr": "KP",
        "img": ""
    },
    {
        "name": "Korea (Republic of)",
        "abr": "KR",
        "img": "https://cdn.countryflags.com/thumbs/south-korea/flag-square-250.png"
    },
    {
        "name": "Kuwait",
        "abr": "KW",
        "img": "https://cdn.countryflags.com/thumbs/kuwait/flag-square-250.png"
    },
    {
        "name": "Kyrgyzstan",
        "abr": "KG",
        "img": "https://cdn.countryflags.com/thumbs/kyrgyzstan/flag-square-250.png"
    },
    {
        "name": "Lao People's Democratic Republic",
        "abr": "LA",
        "img": ""
    },
    {
        "name": "Latvia",
        "abr": "LV",
        "img": "https://cdn.countryflags.com/thumbs/latvia/flag-square-250.png"
    },
    {
        "name": "Lebanon",
        "abr": "LB",
        "img": "https://cdn.countryflags.com/thumbs/lebanon/flag-square-250.png"
    },
    {
        "name": "Lesotho",
        "abr": "LS",
        "img": "https://cdn.countryflags.com/thumbs/lesotho/flag-square-250.png"
    },
    {
        "name": "Liberia",
        "abr": "LR",
        "img": "https://cdn.countryflags.com/thumbs/liberia/flag-square-250.png"
    },
    {
        "name": "Libya",
        "abr": "LY",
        "img": "https://cdn.countryflags.com/thumbs/libya/flag-square-250.png"
    },
    {
        "name": "Liechtenstein",
        "abr": "LI",
        "img": "https://cdn.countryflags.com/thumbs/liechtenstein/flag-square-250.png"
    },
    {
        "name": "Lithuania",
        "abr": "LT",
        "img": "https://cdn.countryflags.com/thumbs/lithuania/flag-square-250.png"
    },
    {
        "name": "Luxembourg",
        "abr": "LU",
        "img": "https://cdn.countryflags.com/thumbs/luxembourg/flag-square-250.png"
    },
    {
        "name": "Macao",
        "abr": "MO",
        "img": "https://cdn.countryflags.com/thumbs/macau/flag-square-250.png"
    },
    {
        "name": "Madagascar",
        "abr": "MG",
        "img": "https://cdn.countryflags.com/thumbs/madagascar/flag-square-250.png"
    },
    {
        "name": "Malawi",
        "abr": "MW",
        "img": "https://cdn.countryflags.com/thumbs/malawi/flag-square-250.png"
    },
    {
        "name": "Malaysia",
        "abr": "MY",
        "img": "https://cdn.countryflags.com/thumbs/malaysia/flag-square-250.png"
    },
    {
        "name": "Maldives",
        "abr": "MV",
        "img": ""
    },
    {
        "name": "Mali",
        "abr": "ML",
        "img": "https://cdn.countryflags.com/thumbs/mali/flag-square-250.png"
    },
    {
        "name": "Malta",
        "abr": "MT",
        "img": "https://cdn.countryflags.com/thumbs/malta/flag-square-250.png"
    },
    {
        "name": "Marshall Islands",
        "abr": "MH",
        "img": ""
    },
    {
        "name": "Martinique",
        "abr": "MQ",
        "img": ""
    },
    {
        "name": "Mauritania",
        "abr": "MR",
        "img": "https://cdn.countryflags.com/thumbs/mauritania/flag-square-250.png"
    },
    {
        "name": "Mauritius",
        "abr": "MU",
        "img": "https://cdn.countryflags.com/thumbs/mauritius/flag-square-250.png"
    },
    {
        "name": "Mayotte",
        "abr": "YT",
        "img": ""
    },
    {
        "name": "Mexico",
        "abr": "MX",
        "img": "https://cdn.countryflags.com/thumbs/mexico/flag-square-250.png"
    },
    {
        "name": "Micronesia (Federated States of)",
        "abr": "FM",
        "img": ""
    },
    {
        "name": "Moldova (Republic of)",
        "abr": "MD",
        "img": ""
    },
    {
        "name": "Monaco",
        "abr": "MC",
        "img": "https://cdn.countryflags.com/thumbs/monaco/flag-square-250.png"
    },
    {
        "name": "Mongolia",
        "abr": "MN",
        "img": "https://cdn.countryflags.com/thumbs/mongolia/flag-square-250.png"
    },
    {
        "name": "Montenegro",
        "abr": "ME",
        "img": "https://cdn.countryflags.com/thumbs/montenegro/flag-square-250.png"
    },
    {
        "name": "Montserrat",
        "abr": "MS",
        "img": ""
    },
    {
        "name": "Morocco",
        "abr": "MA",
        "img": "https://cdn.countryflags.com/thumbs/morocco/flag-square-250.png"
    },
    {
        "name": "Mozambique",
        "abr": "MZ",
        "img": "https://cdn.countryflags.com/thumbs/mozambique/flag-square-250.png"
    },
    {
        "name": "Myanmar",
        "abr": "MM",
        "img": "https://cdn.countryflags.com/thumbs/myanmar/flag-square-250.png"
    },
    {
        "name": "Namibia",
        "abr": "NA",
        "img": "https://cdn.countryflags.com/thumbs/namibia/flag-square-250.png"
    },
    {
        "name": "Nauru",
        "abr": "NR",
        "img": "https://cdn.countryflags.com/thumbs/nauru/flag-square-250.png"
    },
    {
        "name": "Nepal",
        "abr": "NP",
        "img": "https://cdn.countryflags.com/thumbs/nepal/flag-square-250.png"
    },
    {
        "name": "Netherlands",
        "abr": "NL",
        "img": "https://cdn.countryflags.com/thumbs/netherlands/flag-square-250.png"
    },
    {
        "name": "New Caledonia",
        "abr": "NC",
        "img": "https://cdn.countryflags.com/thumbs/new-caledonia/flag-square-250.png"
    },
    {
        "name": "New Zealand",
        "abr": "NZ",
        "img": "https://cdn.countryflags.com/thumbs/new-zealand/flag-square-250.png"
    },
    {
        "name": "Nicaragua",
        "abr": "NI",
        "img": "https://cdn.countryflags.com/thumbs/nicaragua/flag-square-250.png"
    },
    {
        "name": "Niger",
        "abr": "NE",
        "img": "https://cdn.countryflags.com/thumbs/niger/flag-square-250.png"
    },
    {
        "name": "Nigeria",
        "abr": "NG",
        "img": "https://cdn.countryflags.com/thumbs/nigeria/flag-square-250.png"
    },
    {
        "name": "Niue",
        "abr": "NU",
        "img": "https://cdn.countryflags.com/thumbs/niue/flag-square-250.png"
    },
    {
        "name": "Norfolk Island",
        "abr": "NF",
        "img": ""
    },
    {
        "name": "North Macedonia",
        "abr": "MK",
        "img": "https://cdn.countryflags.com/thumbs/north-macedonia/flag-square-250.png"
    },
    {
        "name": "Northern Mariana Islands",
        "abr": "MP",
        "img": ""
    },
    {
        "name": "Norway",
        "abr": "NO",
        "img": "https://cdn.countryflags.com/thumbs/norway/flag-square-250.png"
    },
    {
        "name": "Oman",
        "abr": "OM",
        "img": "https://cdn.countryflags.com/thumbs/oman/flag-square-250.png"
    },
    {
        "name": "Pakistan",
        "abr": "PK",
        "img": "https://cdn.countryflags.com/thumbs/pakistan/flag-square-250.png"
    },
    {
        "name": "Palau",
        "abr": "PW",
        "img": "https://cdn.countryflags.com/thumbs/palau/flag-square-250.png"
    },
    {
        "name": "Palestine, State of",
        "abr": "PS",
        "img": ""
    },
    {
        "name": "Panama",
        "abr": "PA",
        "img": "https://cdn.countryflags.com/thumbs/panama/flag-square-250.png"
    },
    {
        "name": "Papua New Guinea",
        "abr": "PG",
        "img": "https://cdn.countryflags.com/thumbs/papua-new-guinea/flag-square-250.png"
    },
    {
        "name": "Paraguay",
        "abr": "PY",
        "img": "https://cdn.countryflags.com/thumbs/paraguay/flag-square-250.png"
    },
    {
        "name": "Peru",
        "abr": "PE",
        "img": "https://cdn.countryflags.com/thumbs/peru/flag-square-250.png"
    },
    {
        "name": "Philippines",
        "abr": "PH",
        "img": ""
    },
    {
        "name": "Pitcairn",
        "abr": "PN",
        "img": ""
    },
    {
        "name": "Poland",
        "abr": "PL",
        "img": "https://cdn.countryflags.com/thumbs/poland/flag-square-250.png"
    },
    {
        "name": "Portugal",
        "abr": "PT",
        "img": "https://cdn.countryflags.com/thumbs/portugal/flag-square-250.png"
    },
    {
        "name": "Puerto Rico",
        "abr": "PR",
        "img": "https://cdn.countryflags.com/thumbs/puerto-rico/flag-square-250.png"
    },
    {
        "name": "Qatar",
        "abr": "QA",
        "img": "https://cdn.countryflags.com/thumbs/qatar/flag-square-250.png"
    },
    {
        "name": "Réunion",
        "abr": "RE",
        "img": ""
    },
    {
        "name": "Romania",
        "abr": "RO",
        "img": "https://cdn.countryflags.com/thumbs/romania/flag-square-250.png"
    },
    {
        "name": "Russian Federation",
        "abr": "RU",
        "img": "https://cdn.countryflags.com/thumbs/russia/flag-square-250.png"
    },
    {
        "name": "Rwanda",
        "abr": "RW",
        "img": "https://cdn.countryflags.com/thumbs/rwanda/flag-square-250.png"
    },
    {
        "name": "Saint Barthélemy",
        "abr": "BL",
        "img": ""
    },
    {
        "name": "Saint Helena, Ascension and Tristan da Cunha",
        "abr": "SH",
        "img": ""
    },
    {
        "name": "Saint Kitts and Nevis",
        "abr": "KN",
        "img": "https://cdn.countryflags.com/thumbs/saint-kitts-and-nevis/flag-square-250.png"
    },
    {
        "name": "Saint Lucia",
        "abr": "LC",
        "img": "https://cdn.countryflags.com/thumbs/saint-lucia/flag-square-250.png"
    },
    {
        "name": "Saint Martin (French part)",
        "abr": "MF",
        "img": ""
    },
    {
        "name": "Saint Pierre and Miquelon",
        "abr": "PM",
        "img": ""
    },
    {
        "name": "Saint Vincent and the Grenadines",
        "abr": "VC",
        "img": "https://cdn.countryflags.com/thumbs/saint-vincent-and-the-grenadines/flag-square-250.png"
    },
    {
        "name": "Samoa",
        "abr": "WS",
        "img": "https://cdn.countryflags.com/thumbs/samoa/flag-square-250.png"
    },
    {
        "name": "San Marino",
        "abr": "SM",
        "img": "https://cdn.countryflags.com/thumbs/san-marino/flag-square-250.png"
    },
    {
        "name": "Sao Tome and Principe",
        "abr": "ST",
        "img": ""
    },
    {
        "name": "Saudi Arabia",
        "abr": "SA",
        "img": "https://cdn.countryflags.com/thumbs/saudi-arabia/flag-square-250.png"
    },
    {
        "name": "Senegal",
        "abr": "SN",
        "img": "https://cdn.countryflags.com/thumbs/senegal/flag-square-250.png"
    },
    {
        "name": "Serbia",
        "abr": "RS",
        "img": "https://cdn.countryflags.com/thumbs/serbia/flag-square-250.png"
    },
    {
        "name": "Seychelles",
        "abr": "SC",
        "img": ""
    },
    {
        "name": "Sierra Leone",
        "abr": "SL",
        "img": "https://cdn.countryflags.com/thumbs/sierra-leone/flag-square-250.png"
    },
    {
        "name": "Singapore",
        "abr": "SG",
        "img": "https://cdn.countryflags.com/thumbs/singapore/flag-square-250.png"
    },
    {
        "name": "Sint Maarten (Dutch part)",
        "abr": "SX",
        "img": ""
    },
    {
        "name": "Slovakia",
        "abr": "SK",
        "img": "https://cdn.countryflags.com/thumbs/slovakia/flag-square-250.png"
    },
    {
        "name": "Slovenia",
        "abr": "SI",
        "img": "https://cdn.countryflags.com/thumbs/slovenia/flag-square-250.png"
    },
    {
        "name": "Solomon Islands",
        "abr": "SB",
        "img": ""
    },
    {
        "name": "Somalia",
        "abr": "SO",
        "img": "https://cdn.countryflags.com/thumbs/somalia/flag-square-250.png"
    },
    {
        "name": "South Africa",
        "abr": "ZA",
        "img": "https://cdn.countryflags.com/thumbs/south-africa/flag-square-250.png"
    },
    {
        "name": "South Georgia and the South Sandwich Islands",
        "abr": "GS",
        "img": ""
    },
    {
        "name": "South Sudan",
        "abr": "SS",
        "img": "https://cdn.countryflags.com/thumbs/south-sudan/flag-square-250.png"
    },
    {
        "name": "Spain",
        "abr": "ES",
        "img": "https://cdn.countryflags.com/thumbs/spain/flag-square-250.png"
    },
    {
        "name": "Sri Lanka",
        "abr": "LK",
        "img": "https://cdn.countryflags.com/thumbs/sri-lanka/flag-square-250.png"
    },
    {
        "name": "Sudan",
        "abr": "SD",
        "img": "https://cdn.countryflags.com/thumbs/sudan/flag-square-250.png"
    },
    {
        "name": "Suriname",
        "abr": "SR",
        "img": "https://cdn.countryflags.com/thumbs/suriname/flag-square-250.png"
    },
    {
        "name": "Svalbard and Jan Mayen",
        "abr": "SJ",
        "img": ""
    },
    {
        "name": "Sweden",
        "abr": "SE",
        "img": "https://cdn.countryflags.com/thumbs/sweden/flag-square-250.png"
    },
    {
        "name": "Switzerland",
        "abr": "CH",
        "img": "https://cdn.countryflags.com/thumbs/switzerland/flag-square-250.png"
    },
    {
        "name": "Syrian Arab Republic",
        "abr": "SY",
        "img": ""
    },
    {
        "name": "Taiwan, Province of China[note 1]",
        "abr": "TW",
        "img": ""
    },
    {
        "name": "Tajikistan",
        "abr": "TJ",
        "img": "https://cdn.countryflags.com/thumbs/tajikistan/flag-square-250.png"
    },
    {
        "name": "Tanzania, United Republic of",
        "abr": "TZ",
        "img": "https://cdn.countryflags.com/thumbs/tanzania/flag-square-250.png"
    },
    {
        "name": "Thailand",
        "abr": "TH",
        "img": "https://cdn.countryflags.com/thumbs/thailand/flag-square-250.png"
    },
    {
        "name": "Timor-Leste",
        "abr": "TL",
        "img": ""
    },
    {
        "name": "Togo",
        "abr": "TG",
        "img": "https://cdn.countryflags.com/thumbs/togo/flag-square-250.png"
    },
    {
        "name": "Tokelau",
        "abr": "TK",
        "img": ""
    },
    {
        "name": "Tonga",
        "abr": "TO",
        "img": "https://cdn.countryflags.com/thumbs/tonga/flag-square-250.png"
    },
    {
        "name": "Trinidad and Tobago",
        "abr": "TT",
        "img": "https://cdn.countryflags.com/thumbs/trinidad-and-tobago/flag-square-250.png"
    },
    {
        "name": "Tunisia",
        "abr": "TN",
        "img": "https://cdn.countryflags.com/thumbs/tunisia/flag-square-250.png"
    },
    {
        "name": "Türkiye",
        "abr": "TR",
        "img": ""
    },
    {
        "name": "Turkmenistan",
        "abr": "TM",
        "img": "https://cdn.countryflags.com/thumbs/turkmenistan/flag-square-250.png"
    },
    {
        "name": "Turks and Caicos Islands",
        "abr": "TC",
        "img": ""
    },
    {
        "name": "Tuvalu",
        "abr": "TV",
        "img": "https://cdn.countryflags.com/thumbs/tuvalu/flag-square-250.png"
    },
    {
        "name": "Uganda",
        "abr": "UG",
        "img": "https://cdn.countryflags.com/thumbs/uganda/flag-square-250.png"
    },
    {
        "name": "Ukraine",
        "abr": "UA",
        "img": "https://cdn.countryflags.com/thumbs/ukraine/flag-square-250.png"
    },
    {
        "name": "United Arab Emirates",
        "abr": "AE",
        "img": "https://cdn.countryflags.com/thumbs/united-arab-emirates/flag-square-250.png"
    },
    {
        "name": "United Kingdom of Great Britain and Northern Ireland",
        "abr": "GB",
        "img": "https://cdn.countryflags.com/thumbs/united-kingdom/flag-square-250.png"
    },
    {
        "name": "United States of America",
        "abr": "US",
        "img": "https://cdn.countryflags.com/thumbs/united-states-of-america/flag-square-250.png"
    },
    {
        "name": "United States Minor Outlying Islands",
        "abr": "UM",
        "img": ""
    },
    {
        "name": "Uruguay",
        "abr": "UY",
        "img": "https://cdn.countryflags.com/thumbs/uruguay/flag-square-250.png"
    },
    {
        "name": "Uzbekistan",
        "abr": "UZ",
        "img": "https://cdn.countryflags.com/thumbs/uzbekistan/flag-square-250.png"
    },
    {
        "name": "Vanuatu",
        "abr": "VU",
        "img": "https://cdn.countryflags.com/thumbs/vanuatu/flag-square-250.png"
    },
    {
        "name": "Venezuela (Bolivarian Republic of)",
        "abr": "VE",
        "img": "https://cdn.countryflags.com/thumbs/venezuela/flag-square-250.png"
    },
    {
        "name": "Viet Nam",
        "abr": "VN",
        "img": "https://cdn.countryflags.com/thumbs/vietnam/flag-square-250.png"
    },
    {
        "name": "Virgin Islands (British)",
        "abr": "VG",
        "img": ""
    },
    {
        "name": "Virgin Islands (U.S.)",
        "abr": "VI",
        "img": ""
    },
    {
        "name": "Wallis and Futuna",
        "abr": "WF",
        "img": ""
    },
    {
        "name": "Western Sahara",
        "abr": "EH",
        "img": ""
    },
    {
        "name": "Yemen",
        "abr": "YE",
        "img": "https://cdn.countryflags.com/thumbs/yemen/flag-square-250.png"
    },
    {
        "name": "Zambia",
        "abr": "ZM",
        "img": "https://cdn.countryflags.com/thumbs/zambia/flag-square-250.png"
    },
    {
        "name": "Zimbabwe",
        "abr": "ZW",
        "img": "https://cdn.countryflags.com/thumbs/zimbabwe/flag-square-250.png"
    }
]`

var countryIconMap = map[string]string{}

func init() {
	var countryIcons = []struct {
		Name string `json:"name"`
		Abr  string `json:"abr"`
		Img  string `json:"img"`
	}{}
	json.Unmarshal([]byte(conutryJson), &countryIcons)
	for _, country := range countryIcons {
		countryIconMap[country.Abr] = country.Img
	}
}

func GetIcon(Abr string) string {
	if icon, ok := countryIconMap[Abr]; ok {
		return icon
	}
	return ""
}
