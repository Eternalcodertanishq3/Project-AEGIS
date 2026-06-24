import { useEffect, useRef } from 'react'
import maplibregl from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css'
import * as pmtiles from 'pmtiles'

interface MapViewerProps {
  pmtilesUrl: string
}

export function MapViewer({ pmtilesUrl }: MapViewerProps) {
  const mapContainer = useRef<HTMLDivElement>(null)
  const map = useRef<maplibregl.Map | null>(null)

  useEffect(() => {
    if (!mapContainer.current) return

    // Register PMTiles protocol if not already registered
    const protocol = new pmtiles.Protocol()
    maplibregl.addProtocol('pmtiles', protocol.tile)

    // Ensure the PMTiles source URL is absolute
    const url = new URL(pmtilesUrl, window.location.origin).href

    map.current = new maplibregl.Map({
      container: mapContainer.current,
      style: {
        version: 8,
        sources: {
          'offline-tiles': {
            type: 'vector',
            url: `pmtiles://${url}`
          }
        },
        layers: [
          {
            id: 'background',
            type: 'background',
            paint: {
              'background-color': '#0f172a'
            }
          },
          {
            id: 'water',
            type: 'fill',
            source: 'offline-tiles',
            'source-layer': 'water',
            paint: {
              'fill-color': '#1e293b'
            }
          },
          {
            id: 'buildings',
            type: 'fill',
            source: 'offline-tiles',
            'source-layer': 'building',
            paint: {
              'fill-color': '#334155',
              'fill-opacity': 0.8
            }
          },
          {
            id: 'roads',
            type: 'line',
            source: 'offline-tiles',
            'source-layer': 'transportation',
            paint: {
              'line-color': '#475569',
              'line-width': 1.5
            }
          }
        ]
      },
      center: [174.7633, -36.8485], // Defaulting to Auckland NZ since the file is nz-buildings
      zoom: 12
    })

    map.current.addControl(new maplibregl.NavigationControl(), 'top-right')

    return () => {
      if (map.current) {
        map.current.remove()
        map.current = null
      }
    }
  }, [pmtilesUrl])

  return <div ref={mapContainer} className="w-full h-full" />
}
