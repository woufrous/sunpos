package sunpos

import (
	// "fmt"
	"math"
	"time"
)

const (
	Rad              float64 = math.Pi / 180.
	EarthMeanRadius  float64 = 6371.01
	AstronomicalUnit float64 = 149597890
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type SunCoordinates struct {
	Azimuth     float64
	ZenithAngle float64
}

func Sunpos(udtTime time.Time, udtLocation Location) *SunCoordinates {
	var (
		// Main variables
		dElapsedJulianDays float64
		dDecimalHours      float64
		dEclipticLongitude float64
		dEclipticObliquity float64
		dRightAscension    float64
		dDeclination       float64

		// Auxiliary variables
		dY float64
		dX float64

		// Return variable
		udtSunCoordinates *SunCoordinates = &SunCoordinates{0, 0}
	)

	// Calculate difference in days between the current Julian Day
	// and JD 2451545.0, which is noon 1 January 2000 Universal Time
	{
		var (
			dJulianDate float64
			liAux1      int
			liAux2      int
		)
		// Calculate time of the day in UT decimal hours
		dDecimalHours = float64(udtTime.Hour()) + (float64(udtTime.Minute())+float64(udtTime.Second())/60.0)/60.0
		// Calculate current Julian Day
		liAux1 = (int(udtTime.Month()) - 14) / 12
		liAux2 = (1461*(udtTime.Year()+4800+liAux1))/4 + (367*(int(udtTime.Month())-2-12*liAux1))/12 - (3*((udtTime.Year()+4900+liAux1)/100))/4 + udtTime.Day() - 32075
		dJulianDate = (float64)(liAux2) - 0.5 + dDecimalHours/24.0
		// Calculate difference between current Julian Day and JD 2451545.0
		dElapsedJulianDays = dJulianDate - 2451545.0
	}

	// Calculate ecliptic coordinates (ecliptic longitude and obliquity of the
	// ecliptic in radians but without limiting the angle to be less than 2*Pi
	// (i.e., the result may be greater than 2*Pi)
	{
		var (
			dMeanLongitude float64
			dMeanAnomaly   float64
			dOmega         float64
		)
		dOmega = 2.1429 - 0.0010394594*dElapsedJulianDays
		dMeanLongitude = 4.8950630 + 0.017202791698*dElapsedJulianDays // Radians
		dMeanAnomaly = 6.2400600 + 0.0172019699*dElapsedJulianDays
		dEclipticLongitude = dMeanLongitude + 0.03341607*math.Sin(dMeanAnomaly) + 0.00034894*math.Sin(2*dMeanAnomaly) - 0.0001134 - 0.0000203*math.Sin(dOmega)
		dEclipticObliquity = 0.4090928 - 6.2140e-9*dElapsedJulianDays + 0.0000396*math.Cos(dOmega)
	}

	// Calculate celestial coordinates ( right ascension and declination ) in radians
	// but without limiting the angle to be less than 2*Pi (i.e., the result may be
	// greater than 2*Pi)
	{
		var dSin_EclipticLongitude float64
		dSin_EclipticLongitude = math.Sin(dEclipticLongitude)
		dY = math.Cos(dEclipticObliquity) * dSin_EclipticLongitude
		dX = math.Cos(dEclipticLongitude)
		dRightAscension = math.Atan2(dY, dX)
		if dRightAscension < 0.0 {
			dRightAscension = dRightAscension + 2*math.Pi
		}
		dDeclination = math.Asin(math.Sin(dEclipticObliquity) * dSin_EclipticLongitude)
	}

	// Calculate local coordinates ( azimuth and zenith angle ) in degrees
	{
		var (
			dGreenwichMeanSiderealTime float64
			dLocalMeanSiderealTime     float64
			dLatitudeInRadians         float64
			dHourAngle                 float64
			dCos_Latitude              float64
			dSin_Latitude              float64
			dCos_HourAngle             float64
			dParallax                  float64
		)
		dGreenwichMeanSiderealTime = 6.6974243242 + 0.0657098283*dElapsedJulianDays + dDecimalHours
		dLocalMeanSiderealTime = (dGreenwichMeanSiderealTime*15 + udtLocation.Longitude) * Rad
		dHourAngle = dLocalMeanSiderealTime - dRightAscension
		dLatitudeInRadians = udtLocation.Latitude * Rad
		dCos_Latitude = math.Cos(dLatitudeInRadians)
		dSin_Latitude = math.Sin(dLatitudeInRadians)
		dCos_HourAngle = math.Cos(dHourAngle)
		udtSunCoordinates.ZenithAngle = (math.Acos(dCos_Latitude*dCos_HourAngle*math.Cos(dDeclination) + math.Sin(dDeclination)*dSin_Latitude))
		dY = -math.Sin(dHourAngle)
		dX = math.Tan(dDeclination)*dCos_Latitude - dSin_Latitude*dCos_HourAngle
		udtSunCoordinates.Azimuth = math.Atan2(dY, dX)
		if udtSunCoordinates.Azimuth < 0.0 {
			udtSunCoordinates.Azimuth = udtSunCoordinates.Azimuth + 2*math.Pi
		}
		udtSunCoordinates.Azimuth = udtSunCoordinates.Azimuth / Rad
		// Parallax Correction
		dParallax = (EarthMeanRadius / AstronomicalUnit) * math.Sin(udtSunCoordinates.ZenithAngle)
		udtSunCoordinates.ZenithAngle = (udtSunCoordinates.ZenithAngle + dParallax) / Rad
	}
	return udtSunCoordinates
}

/*
func main() {
	dt := time.Date(2015, 9, 18, 12, 0, 0, 0, time.UTC)
	sc := sunpos(dt, Location{48.148, -11.573})

	fmt.Printf("Sun pos: %f, %f\n", sc.Azimuth, sc.ZenithAngle)
}
*/
