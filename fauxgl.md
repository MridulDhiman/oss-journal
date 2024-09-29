### Abstract

fauxgl is a CPU based 3D rendering library written in go. I am thinking of exploring it, and understanding how 3D rendering works. Let's see how it goes.
I always procrastinate trying to explore new repositories. But, not this time. 
I am writing this blog to make myself accountable.

I know you might be wondering, why do I want to learn how 3D rendering works. Well even I don't know, what I want to do with this. I don't even have an interest in this shit. But, as we have started it, I will see it through the end. Let's see how it goes. 

So, whoever is reading this, this personal journal might be eternities long. 

### Day 1 (29th September, 2024): 
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
It is represented by Min, Max Vectors in `BoundaryBox` struct.
It has various purposes in 3d graphics and game development:
1. To check whether 3d objects are intersecting/touching or not. It can also be used for collision detection in 3d models. Basically, what we do is we create bounding boxes of each 3D object, and check whether they intersect or not. If they do not intersect, objects have not collided. Implementation: `Intersects()` method.
2. Using this, we can divide our 3D space into multiple parts, and optimize the search by reducing the search space to a specific bounding box. 
Implementation: `Contains(), ContainsBox(), Intersection()` methods are used to determine which partition an object belongs to.
3. Frustum Culling: it is used to determine which objects are visible from camera's point of view. We create a frustum from camera's point of view. Then, the objects whose bounding boxes are completely outside the frustum, can be discarded from the rendering pipeline. 
Implementation: we can use `Intersects()` method to determine if object is in camera point of view or not, by checking if it intersects with the view frustum of the camera.


