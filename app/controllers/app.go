package controllers

import "github.com/revel/revel"
import "net/http"
import "fmt"
import "io/ioutil"
import "encoding/json"

type App struct {
	*revel.Controller
}

type NasaObj struct{
    Name   string
    Url    string
    Danger  bool
    Magnitude   float64
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Nasa(startdate string) revel.Result{
    c.Validation.Required(startdate).Message("startdate is required")
    if c.Validation.HasErrors(){
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(App.Index)
    }
    resp, err := http.Get("https://api.nasa.gov/neo/rest/v1/feed?start_date="+startdate+"&end_date="+startdate+"&api_key=fEEk9AyeFfKTSK387wh9Jhe4GD8B6HNPtUR0b4ma")
    if err != nil {
        fmt.Println("Error while fetching")
        fmt.Println(err) 
    }
    body, err := ioutil.ReadAll(resp.Body)
    if (err != nil) {
        fmt.Println("Error while Reading the body")
        fmt.Println(err)
    }
    var f map[string]interface{}
    errt := json.Unmarshal(body, &f)
    if (errt != nil){
        fmt.Println("Error while unmarshaling")
        fmt.Println(errt)
    }
    var specific_data = f["near_earth_objects"].(map[string]interface{})[startdate].([]interface{})
    var objs []NasaObj
    for i := 0; i < len(specific_data); i++ {
        var nodedata = specific_data[i].(map[string]interface{})
        obj := NasaObj{Name:nodedata["name"].(string),
                        Url:nodedata["nasa_jpl_url"].(string),
                        Danger:nodedata["is_potentially_hazardous_asteroid"].(bool),
                        Magnitude:nodedata["absolute_magnitude_h"].(float64)}
        objs = append(objs, obj)
    }
    defer resp.Body.Close()
    return c.Render(startdate, objs)
}
