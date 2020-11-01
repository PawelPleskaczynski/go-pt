# go-pt

This project is a Monte Carlo path tracer written in Golang that runs on CPU.

![Monkeys](https://i.imgur.com/t7qnzOA.png)

## Features
### Implemented
- Parallel processing on multiple CPU cores
- BVH trees for accelerating intersection tests
- Positionable camera with adjustable focal length and aperture
- Transformations (translation, rotation)
- Materials:
    - universal material with adjustable properties:
        - albedo:
            - texture or color
        - roughness (GGX microfacet model)
        - index of refraction
        - amount of clearcoat
        - roughness of clearcoat
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
- Nishita sky model with a sun
### To-do
- Building scenes from files (probably JSON?)
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

## Example renders
Some of the models downloaded from Morgan McGuire's [Computer Graphics Archive](https://casual-effects.com/data).
The Go gopher was designed by Renee French. (http://reneefrench.blogspot.com/).
The gopher 3D model was made by Takuya Ueda (https://twitter.com/tenntenn).
HDRI image used in one of the examples was downloaded from [HDRI Haven](https://hdrihaven.com/hdri/?h=river_walk_1).
Textures used in one of the examples were downloaded from [TextureCan](https://www.texturecan.com/) and [Texture Haven](https://texturehaven.com/).

Mori knobs - rough copper, rough glass, glossy plastic

![Mori knobs](https://i.imgur.com/ZKOIb0V.png)

A scene with Nishita sky model with sun (denoised in Blender)

![Scene with sky](https://i.imgur.com/4IRl0bB.png)

Render of a scene with materials with different level on roughness from my path tracer on the left and reference render from Blender on the right

![Comparison](https://i.imgur.com/vzurrgh.png)

Four spheres with different textures and materials with HDRI environment map

![HDRI example](https://i.imgur.com/RAUVu5I.png)

Cornell box with mirror box

![Cornell box](https://i.imgur.com/aolJP3j.png)