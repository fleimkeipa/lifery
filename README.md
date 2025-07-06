Lifery is a full-stack web platform that allows users to record, organize, and share life memories and milestones in a rich multimedia diary format. Users can attach photos, videos, audio, and text to each memory and organize them chronologically on a timeline, grouped by life periods (e.g., “University Years,” “First Job”). Each entry supports granular visibility settings (public, friends-only, or private), providing users with full control over their digital narrative.

From a technical perspective, Lifery is built with Go (Golang) on the backend, chosen for its performance, concurrency support, and suitability for scalable microservices. The backend follows Clean Architecture principles to ensure maintainability, modularity, and testability, with a clear separation of concerns across handler, use case, and repository layers. PostgreSQL is used as the relational database for its robustness, extensibility, and ACID compliance.

The frontend is implemented using Vue.js and Nuxt.js, leveraging server-side rendering (SSR) for improved SEO and performance. The interface emphasizes clarity and usability, featuring timeline-based visualization and dynamic form components for multimedia content. For deployment, the project utilizes the Railway cloud platform, integrating with GitHub for automated CI/CD, ensuring reliable delivery and fast iteration cycles.

Overall, Lifery combines thoughtful user experience design with modern, scalable software architecture to provide a private and meaningful space for storytelling and memory archiving—bridging the gap between personal journaling and social sharing.

<img width="466" alt="Ekran Resmi 2025-07-06 19 10 17" src="https://github.com/user-attachments/assets/4ee1a138-0c73-4cf3-8180-93db9b412401" />

Swagger docs: https://lifery-production.up.railway.app/swagger/index.html#/
