## Abstract

fauxgl is a CPU based 3D rendering library written in go. It is a software based 3d rendering library and does not uses any hardware accelerators like GPUS etc.  I am thinking of exploring it, and understanding how 3D rendering works. Let's see how it goes.
I always procrastinate trying to explore new repositories. But, not this time. 
I am writing this blog to make myself accountable.

I know you might be wondering, why do I want to learn how 3D rendering works. Well even I don't know, what I want to do with this. I don't even have an interest in this shit. But, as we have started it, I will see it through the end. Let's see how it goes. 

So, whoever is reading this, this personal journal might be eternities long. 

## Index

 - [Initial Learnings](#initial-learnings)

### Initial Learnings 
I have cloned this repository and started exploring the components. Basically, I can see a lot of mathematical things like vertex, vector, scaler, meshes etc. Bruh, I'm so dead.

**Vector**: 
It is represented by X, Y, Z float64 coordinates. You can perform arithmetic operations on two vectors like add, multiply, subtract, divide, etc., dot and cross product, and many more things.

**Line**: 
formed by 2 vertexes. 

**Mesh**: 
It is used to represent 3D mesh, which consists of triangles(made up of 3 vertices) and lines(made up of 2 vertices). Here, triangle are basic building blocks of the mesh structure, and lines are use for additional features. Each mesh has a bounding box associated with it.

*Diff. Types of Meshes*: 
1. Default Mesh: with both lines and triangles
2. Empty Mesh: without any lines and triangles, just an empty struct
3. Line Mesh: with lines only
4. Triangle Mesh: with triangles only

*Diff. operations on meshes*:
1. `Copy()`: perform deep copy of meshes
2. `Add()`: combines 2 meshes
3. `SetColor()`: Sets the color of all triangles of the mesh
4. `Volume() and SurfaceArea()` to calculate volume and surface area of the mesh.

and lot of other things like transformations, smoothing of meshes, mesh modifications, loading meshes to STL files etc.

**Bounding Boxes**: 
It is a box which signifies the boundary of an object in the 3d space. It is represented by Min, Max Vectors in `BoundaryBox` struct. Each object(whether it be line/sphere/triangle) has a boundary box associated with it in the 3D space.
It has various purposes in 3d graphics and game development:
1. To check whether 3d objects are intersecting/touching or not. It can also be used for collision detection in 3d models. Basically, what we do is we create bounding boxes of each 3D object, and check whether they intersect or not. If they do not intersect, objects have not collided. Implementation: `Intersects()` method.
2. Using this, we can divide our 3D space into multiple parts, and optimize the search by reducing the search space to a specific bounding box. 
Implementation: `Contains(), ContainsBox(), Intersection()` methods are used to determine which partition an object belongs to.
3. Frustum Culling: it is used to determine which objects are visible from camera's point of view. We create a frustum from camera's point of view. Then, the objects whose bounding boxes are completely outside the frustum, can be discarded from the rendering pipeline. 
Implementation: we can use `Intersects()` method to determine if object is in camera point of view or not, by checking if it intersects with the view frustum of the camera.

**Matrix**:
4X4 matrix representation which is standard for 3D transformations. It is used for various transformations like translation (changing the position of object), rotation (of objects around specific axis, like aligning the object to face the camera), scaling (changing the size of objects), projection etc. 
We can create projection matrices like frustum projection matrix (for camera's point of view of objects), often used in VRs.

In a typical rendering pipeline: 
1. object matrix: position and orient 3d models in the scene.
2. view matrix: used to position the camera
3. projection matrix: used to create a desired perspective, like to position objects according to the camera's perspective. 
4. these matrices are then combined and vertices of 3d models are transformed.
5. resulting coordinates are then mapped onto screen space for rendering.

**Rasterization**:
process of converting geometric objects like triangle and lines to pixels on the screen. So, basically we can convert 3d mesh to image using rasterization.

How it works in typical rendering pipeline: 
1. Vertex processing: vertices are transformed from local space to screen space using projections and perspectives according to the camera's point of view.
2. Rasterization: `DrawTriangles()` method used for rasterize 3D mesh with triangles. Vertices of triangle are computed for their position in the screen, `drawClippedTriangle()` clips the triangle according to view frustum.
3. Pixel Color Calculation: for each pixel covered by triangle, color is decided by invoking the fragment shader. Resulting color is blended with existing colors in the `ColorBuffer` for that pixel.
4. Depth Management: `DepthBuffer` stores the depth of each pixel in the screen, like they would be oriented according to camera's point of view, and how it's positioned.

**Clipping**
Clipping triangles and lines against a particular clipping plane, s.t. only certain parts are visible according to view frustum. Thus, it reduces load in subsequent stages of the rendering pipeline, thus reducing rendering times.
1. `ClipPlanes` slice has 6 `ClipPlane` struct with P being point on the plane and N being normal vector on the plane. Each plane is used to trim triangle and lines, according to view frustum of the camera.

**Sutherland - Hodgman Algorithm**:
algorithm to clip polygon (triangle in this case), against multiple clip planes. It finds out which input points are remaining after clipping the shapes, and use these output points in the next stage of rendering pipeline. 