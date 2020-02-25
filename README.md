# go-pt

This project is a Monte Carlo path tracer written in Golang that runs on CPU.
![angel](https://i.imgur.com/HOBWasb.png)

## Features
### Implemented
- Parallel processing on multiple CPU cores
- BVH trees for optimized ray-triangle intersection tests
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
- Partial support for OBJ files (the program can parse triangles and their normals, but for now there's no support for textures or different materials for different parts of the model)
- Normal smoothing
- Textures
    - Generated textures
    - Image textures
    - For now they only work with spheres
### To-do
- Building scenes from files (probably JSON?)
- Transformations (translation, rotation, etc.)
- More primitives and BVH trees for them
    - Constructive solid geometry
- Full support for OBJ files
- Normal maps
- Volumetric rendering
- Importance sampling
- Spectral rendering

## Example renders
Models downloaded from Morgan McGuire's [Computer Graphics Archive](https://casual-effects.com/data).
The Go gopher was designed by Renee French. (http://reneefrench.blogspot.com/) The gopher 3D model was made by Takuya Ueda (https://twitter.com/tenntenn).


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