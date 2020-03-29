# go-pt

This project is a Monte Carlo path tracer written in Golang that runs on CPU.

![angel](https://i.imgur.com/HOBWasb.png)

## Features
### Implemented
- Parallel processing on multiple CPU cores
- BVH trees for accelerating intersection tests
- Positionable camera with adjustable field of view and aperture
- Materials:
    - BSDF material with adjustable properties:
        - albedo:
            - texture or color
        - roughness:
        - index of refraction
        - specularity
        - metalicity
        - transmission
    - Emission material:
        - emission color
- Support for OBJ files:
    - loading vertices, texture coordinates and normals
    - triangle fan triangulation of polygons
    - support for materials from MTL files
    - support for image textures
    - normal smoothing
- Textures
    - Generated textures:
        - checkerboard (based on UVs or coordinates)
        - grid with variable line thickness (based on UVs or coordinates)
    - Image textures
- Environment textures
    - Can be loaded from normal image files or from Radiance HDR files (loaded using [hdr](https://github.com/mdouchement/hdr) library)
### To-do
- Building scenes from files (probably JSON?)
- Transformations (translation, rotation, etc.)
- More primitives and BVH trees for them
    - Constructive solid geometry
- Volumetric rendering
- Importance sampling
- Spectral rendering

## Usage
For now, scene has to be set in `main.go` file, I'm planning to add support for reading scenes from files in the future.
This program has only one external dependency, it's [hdr](https://github.com/mdouchement/hdr) library, to install it, run the following command:

```
go get github.com/mdouchement/hdr
```

If you're not planning to use HDRI environment maps, remove the following line from the top of `main.go` file:

```
_ "github.com/mdouchement/hdr/codec/rgbe"
```

To run the program, type the following command:

```
go run .
```

The program should give you an output, like:

```
2020/03/29 16:42:20 Loading scene...
2020/03/29 16:42:20 Loading cornellbox_objects.obj...
2020/03/29 16:42:31 Loading cornellbox_lights.obj...
2020/03/29 16:42:31 Loading cornellbox_floor.obj...
2020/03/29 16:42:31 Building BVHs...
100.00% (2010603/2010603 triangles, 25/25 objects)
2020/03/29 16:42:44 Built BVHs
2020/03/29 16:42:44 Rendering 25 objects (2010603 triangles) and 0 spheres at 256x256 at 128 samples on 16 cores
100.00% ( 128/ 128)      6.0952337s/frame,        93.005µs sample time, ETA:              0s
2020/03/29 16:43:37 Rendering took 52.3339386s
Average frame time: 6.49947957s, average sample time: 99.173µs
Saving...
```

## Example renders
Some of the models downloaded from Morgan McGuire's [Computer Graphics Archive](https://casual-effects.com/data).
The Go gopher was designed by Renee French. (http://reneefrench.blogspot.com/).
The gopher 3D model was made by Takuya Ueda (https://twitter.com/tenntenn).
HDRI image used in one of the examples was downloaded from [HDRI Haven](https://hdrihaven.com/hdri/?h=river_walk_1).
Textures used in one of the examples were downloaded from [TextureCan](https://www.texturecan.com/) and [Texture Haven](https://texturehaven.com/).

Four spheres with different textures and materials with HDRI environment map

![HDRI example](https://i.imgur.com/RAUVu5I.png)

Mori knob, with rough glass material outside, and rough metallic material with texture inside

![Mori knob](https://i.imgur.com/WSRmlmt.png)

Cornell box

![Cornell box](https://i.imgur.com/dSaLhwd.png)

Mori knob

![Mori knob](https://i.imgur.com/jE4yPNP.png)

Stanford dragon with glossy, glass and metal spheres

![dragon](https://i.imgur.com/iLplu0d.png)

Cornell box with various objects

![box with objects](https://i.imgur.com/V7AuTSD.png)

Sphere with UV texture

![Sphere with UV texture](https://i.imgur.com/ZQDCjSn.png)